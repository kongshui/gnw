package nodeservice

import (
	conf "github.com/kongshui/danmu/conf/node"

	"context"
	"fmt"
	"net"

	"github.com/kongshui/gnw/msg"
)

func Listen(ctx context.Context) {
	// Listen listen a port and accept the conn, then add the conn to client list and start receive message.
	l, err := net.Listen("tcp", conf.GetConf().Server.Addr+":"+conf.GetConf().Server.Port)
	if err != nil {
		panic(err)
	}
	fmt.Println("listen on " + conf.GetConf().Server.Addr + ":" + conf.GetConf().Server.Port)
	for {
		conn, err := l.Accept()
		if err != nil {
			ziLog.Error(fmt.Sprintf("accept connection error: %v", err), debug)
		}
		ziLog.Info(fmt.Sprintf("accept a new connection: %s", conn.RemoteAddr().String()), debug)
		nMsg := msg.NewMsgConn(&conn, true, false)
		go nodeInitMsgConn(ctx, nMsg)
	}
}

// func WebListen() { // 等待etcd连接
// 	r := gin.Default()
// 	r.POST("/ws", message.WebGetFowardHandler)
// 	r.Run(":" + conf.GetConf().Server.WebPort)
// }
