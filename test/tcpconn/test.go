package main

import (
	"fmt"
	"net"
	"sync"

	msgnew "github.com/kongshui/gnw/msg"
	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

var (
	msgPool sync.Pool = sync.Pool{
		New: func() any {
			msg := msgnew.NewMsgConn(nil, true, false)
			return msg
		},
	}
)

func main() {
	msg := msgPool.Get().(msginterface.MsgConn)
	defer msgPool.Put(msg)
	if msg == nil || msg.GetId() != 1 {
		fmt.Println("msg is nil", msg)
	}
}

func TcpConn() {
	l, _ := net.Listen("tcp", "127.0.0.1:8080")
	conn, _ := l.Accept()
	fmt.Println("LocalAddr", conn.LocalAddr().String())
	fmt.Println("RemoteAddr", conn.RemoteAddr().String())
	fmt.Println("LocalAddr network", conn.LocalAddr().Network())
	fmt.Println("RemoteAddr network", conn.RemoteAddr().Network())
	conn.Close()
	l.Close()
}
