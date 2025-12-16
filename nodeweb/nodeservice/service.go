package nodeservice

import (
	"github.com/google/uuid"
	conf "github.com/kongshui/danmu/conf/nodeweb"
	dao "github.com/kongshui/danmu/dao/etcd"
	"github.com/kongshui/danmu/zilog"
)

var (
	config     *conf.Config
	etcdClient *dao.Etcd
	gatewayId  int64
	// first_ctx  = context.Background()
	ziLog    *zilog.LogStruct
	debug    = false
	nodeUuid uuid.UUID
)

func Init(cfg *conf.Config, etcd *dao.Etcd, log *zilog.LogStruct, uid uuid.UUID) {
	config = cfg
	etcdClient = etcd
	ziLog = log
	nodeUuid = uid
	// 初始化消息
	// go getBackDomain(first_ctx)
	// go RegisterToEtcd(first_ctx)
	// Listen(first_ctx)

}
