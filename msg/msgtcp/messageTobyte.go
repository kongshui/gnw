package msg

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// gateway设置写入client结构体
// func ClientMsgContext(msgId uint32, data []byte) (msgMessage []byte) {
// 	msgIdB := make([]byte, 4)
// 	binary.BigEndian.PutUint32(msgIdB, msgId)
// 	msgLength := uint32(len(data))
// 	msgLengthB := make([]byte, 4)
// 	binary.BigEndian.PutUint32(msgLengthB, msgLength)
// 	msgMessage = append(msgMessage, msgIdB...)
// 	msgMessage = append(msgMessage, msgLengthB...)
// 	msgMessage = append(msgMessage, data...)
// 	return
// }

// Nodehe 写入结构体
func MsgContext(uidStr string, msgId pmsg.MessageId, data []byte, extra string) (msgMessage []byte) {
	sData := &pmsg.MessageBody{MsgId: msgId, MessageData: data, Timestamp: time.Now().UnixMilli(), Uuid: uidStr, Extra: extra}
	sendMessage, err := proto.Marshal(sData)
	if err != nil {
		log.Println("NodeMsgContext失败, id:", uidStr, "err:", err, "data:", sData)
	}
	msgLength := uint32(len(sendMessage))
	msgLengthB := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLengthB, msgLength)
	msgMessage = append(msgMessage, msgLengthB...)
	msgMessage = append(msgMessage, sendMessage...)
	return
	// msgIdB := make([]byte, 4)
	// binary.BigEndian.PutUint32(msgIdB, msgId)
	// //msgMessage = append(msgMessage, messageIdB...)
	// switch msgId {
	// case 99999999:
	// 	msgMessage = append(msgMessage, uidByte...)
	// 	msgMessage = append(msgMessage, msgIdB...)
	// 	return
	// default:
	// 	msgLength := uint32(len(data))
	// 	msgLengthB := make([]byte, 4)
	// 	binary.BigEndian.PutUint32(msgLengthB, msgLength)
	// 	msgMessage = append(msgMessage, uidByte...)
	// 	msgMessage = append(msgMessage, msgIdB...)
	// 	msgMessage = append(msgMessage, msgLengthB...)
	// 	msgMessage = append(msgMessage, data...)
	// 	return
	// }
}
