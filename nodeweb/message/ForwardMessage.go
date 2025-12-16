package message

import (
	"fmt"
	"log"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/service"
	"google.golang.org/protobuf/proto"
)

// 从gateway获取forward消息
func fromGatewayGetForwardMessageHandler(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	msg := msgBodyPool.Get().(*pmsg.MessageBody)
	defer func() {
		msg.Reset()
		msgBodyPool.Put(msg)
	}()
	if err := proto.Unmarshal(data, msg); err != nil {
		log.Println("websocket uplink err: ", err)
	}
	msg.Uuid = uidStr
	if err := service.WebsocketMessageFunc(msg); err != nil {
		ziLog.Error(fmt.Sprintf("websocketMessageFunc err,  websocket_uplink_msg: %s errorInfo: %s ", msg.String(), err), debug)
	}
}
