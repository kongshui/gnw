package message

import (
	"fmt"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// 接收网关节点信息
func gatewayInfoMessageGetHandler(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	info := &pmsg.NodeInfoMessage{}
	if err := proto.Unmarshal(data, info); err != nil {
		ziLog.Error(fmt.Sprintf("gateway init json unmarshal error %v", err), debug)
		return
	}
	msgConn.SetUuid(info.GetUuid())
	msgConn.SetName(info.GetName())
	msgConn.SetGroupId(info.GetGroupId())
	msgConn.SetNodeType(int8(info.GetNodeType()))
	msgConn.SetOnline(true)
	MessageMap.AddUid(info.GetUuid(), msgConn)
}

// 发送node节点信息
func NodeInfoMessageSend(msgConn msginterface.MsgConn) error {
	info := &pmsg.NodeInfoMessage{}
	info.Uuid = nodeUuid.String()
	info.NodeType = int32(config.Server.NodeType)
	info.Name = config.Server.Name
	info.GroupId = config.Server.GroupId
	sData, err := proto.Marshal(info)
	if err != nil {
		ziLog.Error(fmt.Sprintf("node init json marshal error %v", err), debug)
		return err
	}
	err = sendMessage("", msgConn, pmsg.MessageId_NodeInfo, sData, "")
	if err != nil {
		ziLog.Error(fmt.Sprintf("NodeInfoMessageSend err:%v", err), debug)
		return fmt.Errorf("NodeInfoMessageSend err:%v", err)
	}
	return nil
}
