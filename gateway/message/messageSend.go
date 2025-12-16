package message

import (
	"fmt"
	"time"

	"github.com/kongshui/gnw/gateway/commconet"
	msginterface "github.com/kongshui/gnw/msg/msginterface"
	msg "github.com/kongshui/gnw/msg/msgtcp"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// tcp Send Message To Node
func sendMessageToNode(msgConn msginterface.MsgConn, msgId pmsg.MessageId, data []byte, extra string) error {
	// // id := uint64(uidByte[0])<<56 | uint64(uidByte[1])<<48 | uint64(uidByte[2])<<40 | uint64(uidByte[3])<<32 | uint64(uidByte[4])<<24 | uint64(uidByte[5])<<16 | uint64(uidByte[6])<<8 | uint64(uidByte[7])
	if msgConn == nil {
		sData := msg.MsgContext("", msgId, data, extra)
		// fmt.Println("sendMessageToNode test", extra)
		err := NodeList.SendMsgToNode("", "", sData, 0)
		if err != nil {
			ziLog.Error(fmt.Sprintf("sendMessageToNode test err:%v", err.Error()), debug)
			return err
		}
		return nil
	}
	sData := msg.MsgContext(msgConn.GetUuid(), msgId, data, extra)
	msgConn.CounterAdd()
	// messageMap.Add(uid, msgConn)
	if err := NodeList.SendMsgToNode(msgConn.GetGroupId(), msgConn.GetName(), sData, 0); err != nil {
		ziLog.Error(fmt.Sprintf("gateway 发送消息至node失败, id:%v err:%v data:%v", msgConn.GetUuid(), err, sData), debug)
		return err
	}
	return nil
}

// tcp send Message to client
func tcpSendMessageToClient(uidStr string, msgId pmsg.MessageId, data []byte, extra string) error {
	// nUuid := uuid.MustParse(uidStr)
	_, err := commconet.MessageMap.GetMsgByUuid(uidStr).MessageWrite(msg.MsgContext(uidStr, msgId, data, extra))
	if err != nil {
		ziLog.Error(fmt.Sprintf("gateway 发送消息至client失败, id:%v err:%v data:%v", uidStr, err, data), debug)
		return err
	}
	return nil
}

// websocket send Message to client
func websocketSendMessageToClient(uidStr string, msgId pmsg.MessageId, data []byte, extra string) error {
	var d []byte

	// uid := uuid.MustParse(uidStr)
	switch msgId {
	case pmsg.MessageId_ForwardAck:
		d = data
	default:
		sData := &pmsg.MessageBody{MsgId: msgId, MessageData: data, Timestamp: time.Now().UnixMilli(), Extra: extra}
		d, _ = proto.Marshal(sData)
	}
	_, err := commconet.MessageMap.GetMsgByUuid(uidStr).MessageWrite(d)
	if err != nil {
		ziLog.Error(fmt.Sprintf("gateway 发送消息至client失败, id:%v err:%v data:%v", uidStr, err, data), debug)
		return err
	}
	return nil
}

// sendMessage to client
func sendMessageToClient(uidStr string, msgId pmsg.MessageId, data []byte, extra string) error {
	switch config.Server.ListenMode {
	case "tcp":
		return tcpSendMessageToClient(uidStr, msgId, data, extra)
	case "websocket":
		return websocketSendMessageToClient(uidStr, msgId, data, extra)
	}
	return nil
}
