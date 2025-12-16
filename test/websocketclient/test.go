package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/kongshui/gnw/msg"
	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	var (
		testCount int32
		count     int32
	)
	handler := map[uint32]func(string, msginterface.MsgConn, []byte, string){
		uint32(pmsg.MessageId_TestMsgAck): func(id string, conn msginterface.MsgConn, data []byte, extra string) {
			atomic.AddInt32(&testCount, 1)
			// strSlices := strings.Split(string(data), "_")
			// com := strSlices[len(strSlices)-1]
			// if string(data) != "test_"+fmt.Sprint(testCount) {
			// 	fmt.Println(time.Now().UnixMilli(), string(data), testCount)
			// 	os.Exit(1)
			// }
			fmt.Println(string(data), testCount, 111111)
		},
		uint32(pmsg.MessageId_ForwardAck): func(id string, conn msginterface.MsgConn, data []byte, extra string) {
			// d := &pmsg.MessageBody{}
			// proto.Unmarshal(data, d)
			testCount++
			fmt.Println(time.Now().UnixMilli(), string(data), 9999999999)
		},
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// u := url.URL{Scheme: "wss", Host: "baishouwan.weifenggame.cn", Path: "/ws"}
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:6667", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	fmt.Println(1111)
	defer conn.Close()
	nMsg := msg.NewMsgConn(conn, false, false)
	nMsg.SetOnline(true)
	go nMsg.ReceiveMessage(handler)
	data := &pmsg.MessageBody{}
	data.MsgId = pmsg.MessageId_TestMsg
	data.Uuid = "test"
	d := &pmsg.MessageBody{}
	d.Uuid = "test"
	d.MsgId = pmsg.MessageId_Forward
	d.Extra = "/websocket_callback"
	for {
		// fmt.Println("count:", count)

		if count-testCount > 1000 {
			break
		}
		atomic.AddInt32(&count, 1)
		time.Sleep(1 * time.Second)
		data.MessageData = []byte("test_" + fmt.Sprint(count))
		dataByte, _ := proto.Marshal(data)
		d.MessageData = dataByte
		dByte, _ := proto.Marshal(d)
		// fmt.Println(time.Now().UnixMilli(), 77777777)
		// var cData WebsocketMessageStruct
		// cData.MessageId = 0
		// cData.MessageType = "test"
		// cData.MessageData = data
		// sdata, _ := json.Marshal(cData)
		// fmt.Println(111111)
		fmt.Println("testCount enter:")
		_, err = nMsg.MessageWrite(dByte)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("count:", count)
		// if count%1000 == 1 {
		// 	time.Sleep(1 * time.Second)
		// }
	}
	// done := make(chan struct{})

	// go func() {
	// 	defer close(done)
	// 	for {
	// 		_, message, err := c.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			return
	// 		}
	// 		log.Printf("recv: %s", message)
	// 	}
	// }()

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case t := <-ticker.C:
	// 		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	// 		if err != nil {
	// 			log.Println("write:", err)
	// 			return
	// 		}
	// 	case <-interrupt:
	// 		log.Println("interrupt")

	// 		// Cleanly close the connection by sending a close message and then
	// 		// waiting (with timeout) for the server to close the connection.
	// 		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 		if err != nil {
	// 			log.Println("write close:", err)
	// 			return
	// 		}
	// 		select {
	// 		case <-done:
	// 		case <-time.After(time.Second):
	// 		}
	// 		return
	// 	}
	// }
}

type (
	WebsocketMessageStruct struct {
		MessageId   uint32 // 消息ID
		MessageType string // 消息类型
		MessageData []byte // 消息数据
	}
)
