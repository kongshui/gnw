package service

import (
	"context"

	"github.com/kongshui/gnw/gateway/commconet"
	"github.com/kongshui/gnw/gateway/gwinit"

	conf "github.com/kongshui/danmu/conf/gateway"
)

var (
	etcdClient   = gwinit.Ectd_client
	config       = conf.GetConf()
	nodeInfoList = commconet.NewNodeInfoList()
	debug        = gwinit.Debug
	ziLog        = &gwinit.Zilog
)

func init() {

	// 设置链接模式
	commconet.CommconetInit(config.Gateway.ConnMode, config.Gateway.GroupMode)
	gatewayInit(context.Background())
	// go messageMap.AutoDelete()
}
