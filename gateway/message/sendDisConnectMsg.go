package message

import (
	"fmt"
	"time"

	"github.com/kongshui/gnw/gateway/commconet"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

func sendDisConnectMsg() {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	start := true
	for {
		<-t.C
		if !start {
			continue
		}
		// fmt.Println("开始发送断联消息", start)
		start = false
		// fmt.Println("开始发送断联消息", commconet.DisconnectMap.GetAll())
		// 遍历断联uid
		for k, v := range commconet.DisconnectMap.GetAll() {
			data := &pmsg.Disconnect{}
			data.RoomId = v.RoomId
			data.UserId = v.UserId
			sdata, err := proto.Marshal(data)
			if err != nil {
				ziLog.Error(fmt.Sprintf("SendDisConnectMsg marshal err:%v", err.Error()), debug)
				continue
			}
			dData := &pmsg.MessageBody{}
			dData.MessageData = sdata
			dData.MsgId = pmsg.MessageId_DisConnect
			dData.Uuid = k
			sdata, err = proto.Marshal(dData)
			if err != nil {
				ziLog.Error(fmt.Sprintf("SendDisConnectMsg MessageBody marshal err:%v", err.Error()), debug)
				continue
			}
			// fmt.Println("发送断联消息", pathMap.Get(k))
			if err := sendMessageToNode(nil, pmsg.MessageId_Forward, sdata, pathMap.Get(k)); err != nil {
				ziLog.Error(fmt.Sprintf("SendDisConnectMsg sendMessageToNode err:%v", err.Error()), debug)
				continue
			}
			commconet.DisconnectMap.Remove(k)
			pathMap.Remove(k)
		}
		start = true
	}
}
