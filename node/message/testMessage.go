package message

import (
	"fmt"
	"time"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// 发送测试消息 uidByte为寻找哪一个client，即为client uuid
func NodeMessageHandler(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	fmt.Println(22222222222)
	d := &pmsg.TestMessage{}
	if err := proto.Unmarshal(data, d); err != nil {
		fmt.Println(err, 2222211111)
	}
	fmt.Println("node 收到消息：", d.GetData(), d.GetId(), 2222233333)
	time.Sleep(1 * time.Second)
	fmt.Println(uidStr, 9999999999999999)
	sendMessage(uidStr, msgConn, 2, data, extra)
}
