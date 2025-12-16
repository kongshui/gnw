package message

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/prototext"
)

func WebGetFowardHandler(c *gin.Context) {
	data := bytePool.Get().(*[]byte)
	err := errorPool.Get().(*error)
	defer bytePool.Put(data)
	defer errorPool.Put(err)
	*data, *err = c.GetRawData()
	if *err != nil {
		ziLog.Error(fmt.Sprintf("WebGetFowardHandler GetRawData err:%v", *err), debug)
		c.JSON(200, gin.H{
			"err_code": 1,
			"err_msg":  (*err).Error(),
		})
		return
	}
	uidStrListStr := c.GetHeader("x-client-uuid")
	if uidStrListStr == "" {
		ziLog.Error("WebGetFowardHandler uidStrListStr is null", debug)

		c.JSON(200, gin.H{
			"err_code": 2,
			"err_msg":  "uid is null",
		})
		return
	}
	var uidStrList []string
	if err := json.Unmarshal([]byte(uidStrListStr), &uidStrList); err != nil {
		ziLog.Error(fmt.Sprintf("WebGetFowardHandler json.Unmarshal err: %v ", err), debug)
		c.JSON(200, gin.H{
			"err_code": 3,
			"err_msg":  "WebGetFowardHandler json.Unmarshal err: " + err.Error(),
		})
	}
	// log.Println("WebGetFowardHandler:", uidStrListStr)
	for _, uidStr := range uidStrList {
		if err := sendMessage(uidStr, getGatewayMsgConnByClientUid(uidStr), pmsg.MessageId_ForwardAck, *data, ""); err != nil {
			ziLog.Error(fmt.Sprintf("WebGetFowardHandler 发送消息失败, err:%v", err), debug)
		}
	}

	c.JSON(200, gin.H{
		"err_code": 0,
		"err_msg":  "",
	})
}

// Sse 链接
func SseConn(ctx context.Context, url string) {
	body := map[string]string{
		"uid": nodeUuid.String(),
	}
	bodyType, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyType))
	if err != nil {
		ziLog.Error(fmt.Sprintf("SseConn Post err:%v", err), debug)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ziLog.Error(fmt.Sprintf("SseConn Post err:%v", err), debug)
		return
	}
	scanner := bufio.NewScanner(resp.Body)
	ch := make(chan struct{})
	go sseGetMessage(ctx, scanner, url, ch)
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	count := 0
	for {
		select {
		case <-ctx.Done():
			if SseCancel.Len() == 0 {
				OneGetBackdDomain()
			}
			return
		case <-ch:
			count = 0
		case <-t.C:
			count++
			if count > 3 {
				// 超过50秒没有收到心跳，断开连接
				ziLog.Error("SseConn timeout", debug)
				SseCancel.Delete(url)
				resp.Body.Close()
				if SseCancel.Len() == 0 {
					OneGetBackdDomain()
				}
				return
			}
		}
	}
}

func sseGetMessage(ctx context.Context, scanner *bufio.Scanner, url string, ch chan struct{}) {
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// Empty line indicates the end of an event
			continue
		}
		protoMsg := &pmsg.SseMessage{}
		if err := prototext.Unmarshal([]byte(line), protoMsg); err != nil {
			ziLog.Error(fmt.Sprintf("SseConn Unmarshal err:%v", err), debug)
			continue
		}
		select {
		case <-ctx.Done():
			ziLog.Error(fmt.Sprintf("ctx done,data: %v", protoMsg.String()), debug)
			return
		default:
			switch protoMsg.GetMessageId() {
			case pmsg.MessageId_Ping:
				ch <- struct{}{}
			default:
				count := 0
				for _, uidStr := range protoMsg.GetUidList() {
					count++
					if count > len(protoMsg.GetUidList()) {
						break
					}
					msgConn := getGatewayMsgConnByClientUid(uidStr)
					if msgConn == nil {
						ziLog.Error(fmt.Sprintf("msgConn is nil, uidByte:%v, msgid:%v, data:%v", uidStr, pmsg.MessageId_ForwardAck.String(), protoMsg.String()), debug)
						continue
					}
					if err := sendMessage(uidStr, msgConn, pmsg.MessageId_ForwardAck, protoMsg.GetData(), ""); err != nil {
						ziLog.Error(fmt.Sprintf("WebGetFowardHandler 发送消息失败, err:%v", err), debug)

					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		ziLog.Error(fmt.Sprintf("sseGetMessage scanner err:%v", err), debug)
		SseCancel.Delete(url)
		return
	}
}
