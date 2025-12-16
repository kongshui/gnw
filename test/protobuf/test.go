package main

import (
	"fmt"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

func main() {
	d := &pmsg.MessageBody{MsgId: pmsg.MessageId_TestMsg, MessageData: []byte("test"), Timestamp: 1, Uuid: "", Extra: "test"}

	data, _ := proto.Marshal(d)
	sData := &pmsg.MessageBody{}
	proto.Unmarshal(data, sData)
	sData.Extra = ""
	sData.MessageData = []byte{}
	c, _ := proto.Marshal(sData)
	proto.Unmarshal(c, sData)
	fmt.Println(sData, d)
	if sData.GetExtra() == "" {
		fmt.Println("ok")
	}
}
