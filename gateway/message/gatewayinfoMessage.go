package message

import (
	"fmt"

	msgtcp "github.com/kongshui/gnw/msg/msgtcp"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/gnw/gateway/commconet"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// 发送网关节点信息
func gatewayInfoMessageSendHandler(msgConn msginterface.MsgConn) error {
	info := &pmsg.NodeInfoMessage{}
	info.Uuid = commconet.Uuid.String()
	info.NodeType = int32(config.Server.NodeType)
	info.Name = config.Server.Name
	info.GroupId = config.Server.GroupId
	data, err := proto.Marshal(info)
	if err != nil {
		ziLog.Error(fmt.Sprintf("gateway init json marshal error %v", err), debug)
		return err
	}
	_, err = msgConn.MessageWrite(msgtcp.MsgContext(msgConn.GetUuid(), pmsg.MessageId_NodeInfo, data, ""))
	return err
	// return sendMessageToNode(msgConn, uint32(pmsg.MessageId_NodeInfo), data, "")
}

// 接收node节点信息
func nodeInfoMessageGetHandler(_ string, msgConn msginterface.MsgConn, data []byte, _ string) {
	info := &pmsg.NodeInfoMessage{}
	if err := proto.Unmarshal(data, info); err != nil {
		ziLog.Error(fmt.Sprintf("gateway init json unmarshal error %v", err), debug)
		return
	}
	msgConn.SetUuid(info.GetUuid())
	gatewayInfoMessageSendHandler(msgConn)
	NodeList.AddNode(msgConn)
}
