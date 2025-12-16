package nodeservice

import (
	"github.com/kongshui/gnw/nodeweb/nodeinit"
)

var (
	config     = nodeinit.Config
	etcdClient = nodeinit.Ectd_client
	gatewayId  int64
	// first_ctx  = context.Background()
	ziLog    = &nodeinit.Zilog
	debug    = false
	nodeUuid = nodeinit.NodeUuid
)

func init() {
	// 初始化消息
	// go getBackDomain(first_ctx)
	// go RegisterToEtcd(first_ctx)
	// Listen(first_ctx)
}
