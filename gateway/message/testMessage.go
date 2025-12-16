package message

import (
	"fmt"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

// // gateway接测试消息并发送至node节点
// func testMessageHandler(id uint64, msgConn msginterface.MsgConn, data []byte) {
// 	fmt.Println(11111111111111111)

// 	d := &pmsg.TestMessage{}
// 	if err := proto.Unmarshal(data, d); err != nil {
// 		fmt.Println(err, 1111111)
// 	}
// 	fmt.Println(d.GetData(), d.GetId(), 666666)

// 	sendMessageToNode(msgConn, 2, data)
// }

// websocket
func TestMessageHandler(id string, msgConn msginterface.MsgConn, data []byte, extra string) {
	fmt.Println(11111111111111111)

	// d := &pmsg.TestMessage{}
	// if err := proto.Unmarshal(data, d); err != nil {
	// 	fmt.Println(err, 1111111)
	// }
	// fmt.Println(d.GetData(), d.GetId(), 666666)

	sendMessageToNode(msgConn, 0, data, extra)
}

// gateway 接收node消息并发送给client节点
func TcpMessageHandler(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	fmt.Println(33333333333)
	// d := &pmsg.TestMessage{}
	// if err := proto.Unmarshal(data, d); err != nil {
	// 	fmt.Println(err, 222333333311111)
	// }
	// fmt.Println("gateway 收到消息：", d.GetData(), d.GetId(), 33333332222222)
	// fmt.Println(msgConn.GetGroupId(), msgConn.GetName(), msgConn.GetId(), 333338555555)
	sendMessageToClient(uidStr, 2, data, extra)
}

// websocket
func WebsocketMessageHandler(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	fmt.Println(33333333333)
	sendMessageToClient(uidStr, 2, data, extra)
}
