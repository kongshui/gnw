package msg

import (
	"context"
	"net"
	"reflect"

	msgwb "github.com/kongshui/gnw/msg/msgWebsocket"
	msginterface "github.com/kongshui/gnw/msg/msginterface"
	msgtcp "github.com/kongshui/gnw/msg/msgtcp"

	"github.com/gorilla/websocket"
)

func NewMsgConn(conn any, isTcp bool, isHeartBeatOpen bool) msginterface.MsgConn {
	if conn == nil {
		if isTcp {
			return msgtcp.NewMsgConn(nil, isHeartBeatOpen)
		} else {
			return msgwb.NewMsgConn(nil, isHeartBeatOpen)
		}
	}
	switch reflect.TypeOf(conn).String() {
	case "*net.Conn":
		return msgtcp.NewMsgConn(*conn.(*net.Conn), isHeartBeatOpen)
	case "*websocket.Conn":
		return msgwb.NewMsgConn(conn.(*websocket.Conn), isHeartBeatOpen)
	default:
		return nil
	}
}

// NewMsgConnWithCtx 创建消息连接with ctx 和 conn
func NewMsgConnWithCtx(ctx context.Context, cancel context.CancelFunc, conn any, isHeartBeatOpen bool) msginterface.MsgConn {
	switch reflect.TypeOf(conn).String() {
	case "*net.Conn":
		return msgtcp.NewMsgConnWithCtx(ctx, cancel, *conn.(*net.Conn), isHeartBeatOpen)
	case "*websocket.Conn":
		return msgwb.NewMsgConnWithCtx(ctx, cancel, conn.(*websocket.Conn), isHeartBeatOpen)
	default:
		return nil
	}
}
