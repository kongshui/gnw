package message

import (
	"fmt"
	"path"
	"time"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
	msg "github.com/kongshui/gnw/msg/msgtcp"

	"github.com/kongshui/danmu/common"
	"github.com/kongshui/danmu/model/pmsg"
	"google.golang.org/protobuf/proto"
)

// 发送消息
func sendMessage(uidStr string, msgConn msginterface.MsgConn, msgid pmsg.MessageId, data []byte, extra string) error {
	_, err := msgConn.MessageWrite(msg.MsgContext(uidStr, msgid, data, extra))
	if err != nil {
		ziLog.Error(fmt.Sprintf("node发送消息失败, lableMessageId:%v, msgid:%v, err:%v, data:%v", uidStr, msgid, err, string(data)), debug)

		return err
	}
	return nil
}

// 通过clientuid 查询gateway msgConn
func getGatewayMsgConnByClientUid(clientUid string) msginterface.MsgConn {
	gatewayUid, err := ectdClient.Client.Get(first_ctx, path.Join("/", config.Project, common.Uuid_Register_key, clientUid))
	if err != nil || len(gatewayUid.Kvs) == 0 {
		ziLog.Error(fmt.Sprintf("获取gateway uid失败, err:%v, key:%v", err, common.Uuid_Register_key+clientUid), debug)
		return nil
	}
	// uid := uuid.MustParse(string(gatewayUid.Kvs[0].Value))
	return MessageMap.GetMsgByUuid(string(gatewayUid.Kvs[0].Value))
}

// 发送消息到gateway

func SendMessageToGateway(msgId pmsg.MessageId, uidList []string, data []byte) error {
	dataBody := &pmsg.MessageBody{
		MsgId:       msgId,
		MessageData: data,
		Timestamp:   time.Now().UnixMilli(),
	}
	requestBody, err := proto.Marshal(dataBody)
	if err != nil {
		return fmt.Errorf("SendMessageToGateway proto.Marshal err:%v", err)
	}
	var (
		errList = make([]string, 0)
	)
	for _, uidStr := range uidList {
		msgConn := getGatewayMsgConnByClientUid(uidStr)
		if msgConn == nil {
			errList = append(errList, uidStr)
			continue
		}
		if err := sendMessage(uidStr, msgConn, pmsg.MessageId_ForwardAck, requestBody, ""); err != nil {
			errList = append(errList, uidStr)
			continue
		}
	}
	if len(errList) > 0 {
		return fmt.Errorf("发送消息到gateway失败, errList:%v, msgid:%v, data:%v", errList, msgId, dataBody.String())
	}
	return nil
}
