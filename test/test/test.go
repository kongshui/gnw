package main

import (
	"fmt"
	"net"
	"time"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
	msg "github.com/kongshui/gnw/msg/msgtcp"
	bbb "github.com/kongshui/gnw/test/test/aaa/bbb"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

func main() {
	bbb.BBB()
	handler := map[uint32]func(string, msginterface.MsgConn, []byte, string){
		1: func(id string, conn msginterface.MsgConn, data []byte, extra string) {
			fmt.Println("1")
			fmt.Println(string(data))
		},
		2: func(id string, conn msginterface.MsgConn, data []byte, extra string) {
			fmt.Println(time.Now().UnixMilli(), 9999999999)
			d := &pmsg.MessageBody{}
			proto.Unmarshal(data, d)
			fmt.Println(d.MsgId, d.MessageData, 6666666666666)
		},
	}
	d := &pmsg.MessageBody{}
	d.MessageData = []byte("test")
	d.MsgId = pmsg.MessageId_TestMsg
	count := 0
	conn, _ := net.Dial("tcp", "127.0.0.1:6666")
	defer conn.Close()
	nMsg := msg.NewMsgConn(conn, false)
	nMsg.SetOnline(true)
	go nMsg.ReceiveMessage(handler)
	for {
		count++
		time.Sleep(1 * time.Second)
		d.MessageData = []byte("test" + fmt.Sprint(count))
		data, _ := proto.Marshal(d)
		fmt.Println(time.Now().UnixMilli(), 77777777)
		sdata := msg.MsgContext("", pmsg.MessageId_GetVersionTopHundred, data, "")
		_, err := conn.Write(sdata)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("count:", count)
		if count%10 == 1 {
			time.Sleep(1 * time.Second)
		}
	}

}
