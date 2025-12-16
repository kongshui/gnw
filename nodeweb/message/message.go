package message

import (
	"context"
	"sync"

	"github.com/google/uuid"
	msg "github.com/kongshui/gnw/msg/router"

	conf "github.com/kongshui/danmu/conf/nodeweb"
	dao "github.com/kongshui/danmu/dao/etcd"
	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/zilog"
)

const (
// message_id_timeout = 3600 //消息过期时间
)

var (
	// rdb           = dao_redis.GetRedisClient()
	ectdClient *dao.Etcd
	MessageMap = msg.NewRouterClientMap() //uid和msgConn map
	config     *conf.Config               //配置
	// nats_client    = nats.NatsInit(config.Nats.Addr) //nats
	Handler       = msg.NewMessageHandler() //msgid和handler
	messageIdType = msg.NewIdTypeMap()      //messid和类型map
	first_ctx     = context.Background()    // 初始化ctx
	ziLog         *zilog.LogStruct
	nodeUuid      uuid.UUID
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

func Init(etcd *dao.Etcd, cfg *conf.Config, zlog *zilog.LogStruct, uid uuid.UUID) {
	ectdClient = etcd
	config = cfg
	ziLog = zlog
	nodeUuid = uid
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
