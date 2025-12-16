package msg

import (
	"encoding/binary"
	"io"
	"log"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	pmsg "github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// node消息持续读取数据
func (c *TcpConn) receiveMsg(handler map[uint32]func(string, msginterface.MsgConn, []byte, string)) {
	defer c.Close()
	for {
		msgLenth, err := c.readMsgHeader()
		if err != nil {
			c.SetOnline(false)
			c.Close()
			c.Cancel()
			log.Println("链接断开3")
			log.Println("读取头length错误", err)
			return

		}
		msgData, err := c.readMsgData(msgLenth)
		if err != nil {
			c.SetOnline(false)
			c.Close()
			c.Cancel()
			log.Println("链接断开4")
			log.Println("读取body数据错误", err)
			return
		}
		// c.CounterAdd()
		// log.Println("收到消息: labelMessageId: ", msgData.GetUuid(), "msgId:", msgData.GetMessageId(), "msgData: ", msgData.GetMessageData())
		// 心跳
		switch msgData.GetMsgId() {
		case pmsg.MessageId_Ping:
			c.ReceivePing()
		case pmsg.MessageId_Pong:
			c.ReceivePong()
		default:
			handle, ok := handler[uint32(msgData.GetMsgId())]
			if ok {
				go handle(msgData.GetUuid(), c, msgData.GetMessageData(), msgData.GetExtra())
			}
		}
	}
}

// Node接收数据
func (c *TcpConn) ReceiveMessage(handler map[uint32]func(string, msginterface.MsgConn, []byte, string)) {
	//发送心跳
	go c.keepalived()

	//接收消息
	c.receiveMsg(handler)
	log.Println("退出接收数据:" + c.RemoteAddr().String())
}

// node读取messageId, 暂时弃用
func (c *TcpConn) ReadMsgId() ([]byte, error) {
	// var msgId uint64
	// err := binary.Read(c.conn, binary.BigEndian, &msgId)
	// if err != nil {
	// 	return 0, err
	// }
	// return msgId, nil
	tmp := make([]byte, 16)
	// _, err := conn.Read(tmp), 后期可以测试这两个方法是否有问题
	_, err := io.ReadFull(c.conn, tmp)
	if err != nil {
		log.Println("读取头messageId错误", err)
		return []byte{}, err
	}
	return tmp, err
}

// 读取消息头
func (c *TcpConn) readMsgHeader() (uint32, error) {
	tmp := make([]byte, 4)
	_, err := io.ReadFull(c.conn, tmp)
	if err != nil {
		return 0, err
	}
	heard := binary.BigEndian.Uint32(tmp)
	return heard, nil
}

// 读取body的数据
func (c *TcpConn) readMsgData(length uint32) (*pmsg.MessageBody, error) {
	data := &pmsg.MessageBody{}
	tmp := make([]byte, length)
	// _, err := conn.Read(tmp), 后期可以测试这两个方法是否有问题
	_, err := io.ReadFull(c.conn, tmp)
	if err != nil {
		return data, err
	}
	if err := proto.Unmarshal(tmp, data); err != nil {
		return data, err
	}
	return data, err
}
