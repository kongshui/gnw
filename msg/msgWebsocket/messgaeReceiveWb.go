package msg

import (
	"log"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/gorilla/websocket"
	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// Client接收数据
func (c *WsConn) ReceiveMessage(handler map[uint32]func(string, msginterface.MsgConn, []byte, string)) {
	//发送心跳
	go c.keepalived()

	//接收消息
	go c.receiveMsg(handler)

	<-c.GetCtx().Done()
	log.Println("退出接收数据" + c.RemoteAddr().String())
}

// 接收数据
func (c *WsConn) receiveMsg(handler map[uint32]func(string, msginterface.MsgConn, []byte, string)) {
	defer c.Close()
	for {
		mt, gdata, err := c.GetConn().ReadMessage()
		if err != nil {
			log.Println("接收websocket消息错误 GetConn：", err)
			return
		}
		// // fmt.Println("websocket接收到消息：", string(gdata))
		switch mt {
		case websocket.PingMessage:
			c.ReceivePing()
		case websocket.PongMessage:
			c.ReceivePong()
		default:
			data := &pmsg.MessageBody{}
			err = proto.Unmarshal(gdata, data)
			if err != nil {
				log.Println("接收websocket消息错误：", err, "消息为：", gdata)
				continue
			}
			switch data.GetMsgId() {
			case pmsg.MessageId_Ping:
				c.ReceivePing()
			case pmsg.MessageId_Pong:
				c.ReceivePong()
			default:
				handle, ok := handler[uint32(data.GetMsgId())]
				if ok {
					go handle(data.GetUuid(), c, data.GetMessageData(), data.GetExtra())
				}
			}
		}
	}
}
