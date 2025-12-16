package message

import (
	"github.com/kongshui/gnw/gateway/commconet"
	"github.com/kongshui/gnw/gateway/gwinit"
	msg "github.com/kongshui/gnw/msg/router"

	"github.com/kongshui/danmu/common"

	"github.com/kongshui/danmu/model/pmsg"

	conf "github.com/kongshui/danmu/conf/gateway"
)

var (
	config = conf.GetConf()
	// nodeInfoList  = commconet.NewNodeInfoList()
	NodeList = commconet.NewNodeList()
	// Handler       = msg.New_message_handler()
	messageIdType = msg.NewIdTypeMap()
	Handler       = msg.NewMessageHandler()
	pathMap       = common.NewStringMap()
	ectdClient    = gwinit.Ectd_client
	ziLog         = &gwinit.Zilog
	debug         = gwinit.Debug
)

func init() {
	// 注册消息
	// Handler.Register(0, TestMessageHandler)
	// Handler.Register(1, GetLoadAvg)
	// Handler.Register(2, WebsocketMessageHandler)
	//注册转发信息
	Handler.Register(uint32(pmsg.MessageId_Forward), forwardMessage)
	//注册node过来的转发信息
	Handler.Register(uint32(pmsg.MessageId_ForwardAck), fromNodeGetForwardMessage)
	// 注册nodeinfo信息
	Handler.Register(uint32(pmsg.MessageId_NodeInfo), nodeInfoMessageGetHandler)
	// 注册load消息
	Handler.Register(uint32(pmsg.MessageId_NodeLoad), getLoadAvgMessage)
	//重新登录
	Handler.Register(uint32(pmsg.MessageId_ReLogin), reLoginMessageHandler)

	// 注册id和type
	messageIdType.Register(uint32(pmsg.MessageId_NodeLoad), pmsg.MessageId_NodeLoad.String())
	messageIdType.Register(uint32(pmsg.MessageId_Forward), pmsg.MessageId_Forward.String())
	messageIdType.Register(uint32(pmsg.MessageId_ForwardAck), pmsg.MessageId_ForwardAck.String())
	messageIdType.Register(uint32(pmsg.MessageId_DisConnect), pmsg.MessageId_DisConnect.String())
	go sendDisConnectMsg()
}
