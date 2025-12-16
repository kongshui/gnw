package service

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/kongshui/gnw/gateway/message"

	"github.com/gorilla/websocket"
)

// websocket handler
func webSocket_Handler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ziLog.Error(fmt.Sprintf("webSocket_Handler err:%v", err), debug)
		return
	}
	ziLog.Info(fmt.Sprintf("webSocket connect, ip: %v", conn.RemoteAddr()), debug)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	message.NodeList.GatewayAddClient(ctx, cancel, conn, message.Handler)
	ziLog.Info(fmt.Sprintf("websocket close,client: %v", conn.RemoteAddr()), debug)
}

// websocketCreate create websocket
func websocketCreate() {
	http.HandleFunc("/ws", webSocket_Handler)
	if err := http.ListenAndServe(config.Server.Addr+":"+config.Server.Port, nil); err != nil {
		ziLog.Error(fmt.Sprintf("websocketCreate error:%v", err), debug)
	}
	// if err := http.ListenAndServeTLS(config.Server.Addr+":"+config.Server.Port, config.Server.SslCert, config.Server.SslKey, nil); err != nil {
	// 	log.Println("websocketCreate error:", err)
	// }
}

// tcp handler
func tcpConnListen(ctx context.Context) {
	l, err := net.Listen("tcp", config.Server.Addr+":"+config.Server.Port)
	if err != nil {
		panic(err)
	}
	cCtx, cancel := context.WithCancel(ctx)
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go message.NodeList.GatewayAddClient(cCtx, cancel, &conn, message.Handler)
	}
}

func Listen(ctx context.Context) {
	fmt.Println("listen on "+config.Server.Addr+":"+config.Server.Port, "mode:", config.Server.ListenMode)
	switch config.Server.ListenMode {
	case "tcp":
		tcpConnListen(ctx)
	case "websocket":
		websocketCreate()
	default:
		ziLog.Error(fmt.Sprintf("listen mode error:%v", config.Server.ListenMode), debug)
	}
}
