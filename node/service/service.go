package service

import (
	"context"

	"github.com/kongshui/gnw/node/nodeinit"

	conf "github.com/kongshui/danmu/conf/node"
)

var (
	config     = conf.GetConf()
	etcdClient = nodeinit.Ectd_client
	gatewayId  int64
	first_ctx  = context.Background()
	ziLog      = &nodeinit.Zilog
	debug      = false
	nodeUuid   = nodeinit.NodeUuid
)

func init() {
	go getBackDomain(first_ctx)
}
