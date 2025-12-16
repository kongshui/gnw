package main

import (
	"context"
	"time"

	"github.com/kongshui/gnw/gateway/service"
)

var (
// nodeList = commconet.NewNodeList() //节点列表
// mode     = "hash"                  //hash模式
// is_group = false                   //是否分组
)

func main() {
	var cstZone = time.FixedZone("CST", 8*3600)
	time.Local = cstZone
	ctx := context.Background()
	qCtx, cancle := context.WithCancel(ctx)
	defer cancle()
	// go service.GetNodeList(ctx)

	// go service.GatewayInit(ctx)

	// go Listen(ctx)
	time.Sleep(5 * time.Second)
	service.Listen(qCtx)
}
