package msg

import (
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
)

// 发送心跳
func (c *TcpConn) keepalived() {
	defer c.Close()
	t := time.NewTicker(3 * time.Second)
	bool := true
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if !bool {
				continue
			}

			if err := c.Ping(); err != nil {
				log.Println("心跳发送失败,", "链接为：", c.RemoteAddr(), "本地为：", c.LocalAddr(), "错误为：", err)
				c.Close()
				c.SetState(5)
				c.SetOnline(false)
				c.cancel()
				return
			}
			// log.Println("心跳发送成功")
		case <-c.GetCtx().Done():
			return
		}
	}
}

// ping 心跳
func (c *TcpConn) Ping() error {
	_, err := c.MessageWrite(MsgContext("", pmsg.MessageId_Ping, []byte{}, ""))
	if err != nil {
		return err
	}

	if !c.heartBeatOpen {
		return nil
	}
	atomic.AddInt32(&c.heartBeatCount, 1)

	if c.heartBeatCount > 20 {
		return errors.New("心跳累计次数过多" + c.RemoteAddr().String())
	}
	return nil
}

// pong 心跳回复
func (c *TcpConn) Pong() error {
	_, err := c.MessageWrite(MsgContext("", pmsg.MessageId_Pong, []byte{}, ""))
	if err != nil {
		return err
	}
	return nil
}

// ReceivePong 接收心跳回复
func (c *TcpConn) ReceivePong() error {
	atomic.StoreInt32(&c.heartBeatCount, 0)
	return nil
}

// ReceivePing 接收心跳
func (c *TcpConn) ReceivePing() error {
	// 接收心跳
	return c.Pong()
}
