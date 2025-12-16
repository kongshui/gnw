package nodeservice

import (
	"context"
	"strings"
	"sync/atomic"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
	"github.com/kongshui/gnw/node/message"
)

// nodeInitMsgConn
func nodeInitMsgConn(ctx context.Context, msg msginterface.MsgConn) {
	//设置ctx
	cCtx, cancel := context.WithCancel(ctx)
	msg.SetOnline(true)
	// 初始化信息
	addr, port := strings.Split(msg.RemoteAddr().String(), ":")[0], strings.Split(msg.RemoteAddr().String(), ":")[1]
	msg.SetAddr(addr)
	msg.SetPort(port)
	msg.SetCtx(cCtx)
	msg.SetCancel(cancel)
	atomic.AddInt64(&gatewayId, 1)
	msg.SetId(gatewayId)
	//发送load信息
	message.NodeInfoMessageSend(msg)
	// go message.LoadMessageHandler(msg.GetCtx(), msg)
	// 接收消息
	msg.ReceiveMessage(message.Handler)
}
