package message

import (
	"fmt"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/model/pmsg"
)

// forwardMessage 转发消息
func forwardMessage(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	// log.Println("转发消息", extra)
	pathMap.Add(msgConn.GetUuid(), extra)
	// fmt.Println("转发消息", uidStr)
	if uidStr == "" {
		ziLog.Error("转发消息uidStr为空", debug)
		msgConn.SetOnline(false)
		msgConn.Close()
		msgConn.Cancel()
		return
	}
	// fmt.Println(88888)
	if err := sendMessageToNode(msgConn, pmsg.MessageId_Forward, data, extra); err != nil {
		ziLog.Error(fmt.Sprintf("转发消息toNode失败 %v", err), debug)
		sendMessageToClient(uidStr, pmsg.MessageId_FrontSendMessageError, data, extra)
	}
}

// fromNodeGetForwardMessage 从node获取转发消息
func fromNodeGetForwardMessage(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	// fmt.Println("从node获取转发消息")
	if err := sendMessageToClient(uidStr, pmsg.MessageId_ForwardAck, data, extra); err != nil {
		ziLog.Error(fmt.Sprintf("转发消息toClient失败 %v", err), debug)
	}
}
