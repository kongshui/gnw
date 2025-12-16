package message

import (
	"context"
	"sync"

	msg "github.com/kongshui/gnw/msg/router"
	"github.com/kongshui/gnw/nodeweb/nodeinit"

	"github.com/kongshui/danmu/model/pmsg"
)

const (
// message_id_timeout = 3600 //消息过期时间
)

var (
	// rdb           = dao_redis.GetRedisClient()
	ectdClient = nodeinit.Ectd_client
	MessageMap = msg.NewRouterClientMap() //uid和msgConn map
	config     = nodeinit.Config          //配置
	// nats_client    = nats.NatsInit(config.Nats.Addr) //nats
	Handler       = msg.NewMessageHandler() //msgid和handler
	messageIdType = msg.NewIdTypeMap()      //messid和类型map
	first_ctx     = context.Background()    // 初始化ctx
	ziLog         = nodeinit.Zilog
	nodeUuid      = nodeinit.NodeUuid
	debug         bool // 是否是debug模式
	// bytePool       sync.Pool = sync.Pool{
	// 	New: func() any {
	// 		out := make([]byte, 0)
	// 		return &out
	// 	},
	// }
	// errorPool sync.Pool = sync.Pool{
	// 	New: func() any {
	// 		var err error
	// 		return &err
	// 	},
	// }
	msgBodyPool sync.Pool = sync.Pool{New: func() any {
		return &pmsg.MessageBody{}
	}}
)

func init() {
	go MessageMap.AutoDelete()
	// 发送load
	go loadMessageHandler()
	//
	if config.Logging.Level == "debug" {
		debug = true
	}
	// 注册转发信息
	Handler.Register(uint32(pmsg.MessageId_Forward), fromGatewayGetForwardMessageHandler)
	// 注册gateway发过来的info信息
	Handler.Register(uint32(pmsg.MessageId_NodeInfo), gatewayInfoMessageGetHandler)
	// 注册id和type
	messageIdType.Register(0, "NodeMessageHandler")
	messageIdType.Register(uint32(pmsg.MessageId_NodeInfo), "MessageId_NodeInfo")
	messageIdType.Register(uint32(pmsg.MessageId_Forward), "MessageId_Forward")
}
