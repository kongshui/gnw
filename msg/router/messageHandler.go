package msg

import (
	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

type message_handler map[uint32]func(string, msginterface.MsgConn, []byte, string) //消息处理器
// 消息注册
func (message message_handler) Register(msgId uint32, handler func(string, msginterface.MsgConn, []byte, string)) {
	message[msgId] = handler
}

// new消息处理器
func NewMessageHandler() message_handler {
	return make(map[uint32]func(string, msginterface.MsgConn, []byte, string))
}
