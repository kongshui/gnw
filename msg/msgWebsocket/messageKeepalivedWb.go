package msg

import (
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/kongshui/danmu/model/pmsg"
)

func (c *WsConn) keepalived() {
	if !c.heartBeatOpen {
		return
	}
	defer c.Close()
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-c.GetCtx().Done():
			return
		case <-t.C:
			if err := c.Ping(); err != nil {
				log.Println("心跳发送失败,", "链接为：", c.RemoteAddr(), "本地为：", c.LocalAddr(), "错误为：", err)
				c.Close()
				c.SetState(5)
				c.SetOnline(false)
				c.cancel()
				return
			}
		}
	}
}

// ping 心跳
func (c *WsConn) Ping() error {
	data := &pmsg.MessageBody{MsgId: pmsg.MessageId_Ping, MessageData: []byte{}, Timestamp: time.Now().UnixMilli(), Uuid: "", Extra: ""}
	sendMessage, _ := proto.Marshal(data)
	_, err := c.MessageWrite(sendMessage)
	if err != nil {
		return err
	}
	atomic.AddInt32(&c.heartBeatCount, 1)
	if c.heartBeatCount > 20 {
		return errors.New("心跳累计次数过多" + c.RemoteAddr().String())
	}
	return nil
}

// pong 心跳回复
func (c *WsConn) Pong() error {
	// 发送pong消息
	data := &pmsg.MessageBody{MsgId: pmsg.MessageId_Pong, MessageData: []byte{}, Timestamp: time.Now().UnixMilli(), Uuid: "", Extra: ""}
	sendMessage, _ := proto.Marshal(data)
	_, err := c.MessageWrite(sendMessage)
	if err != nil {
		return err
	}
	return nil
}

// ReceivePong 接收心跳回复
func (c *WsConn) ReceivePong() error {
	atomic.StoreInt32(&c.heartBeatCount, 0)
	return nil
}

// ReceivePing 接收心跳
func (c *WsConn) ReceivePing() error {
	// 接收心跳
	return c.Pong()
}
