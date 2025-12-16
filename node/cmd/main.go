package main

import (
	"context"
	"time"

	"github.com/kongshui/gnw/node/service"
)

func main() {
	var cstZone = time.FixedZone("CST", 8*3600)
	time.Local = cstZone
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx
	// time.Sleep(10 * time.Second)
	go service.Listen(ctx)
	// go service.WebListen()
	service.RegisterToEtcd(ctx)
	// <-ctx.Done()
}
