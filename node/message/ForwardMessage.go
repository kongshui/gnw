package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/common"
	"github.com/kongshui/danmu/model/pmsg"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 从gateway获取forward消息
func fromGatewayGetForwardMessageHandler(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	defer func() {
		if r := recover(); r != nil {
			ziLog.Error(fmt.Sprintf("fromGatewayGetForwardMessageHandler recover: %v", r), debug)
			sendMessage(uidStr, msgConn, pmsg.MessageId_FrontSendMessageError, data, extra)
		}
	}()
	//通过uid获取roomid
	//访问后端web服务器
	header := map[string]string{
		"Content-Type":    "application/json", // client uuid
		"x-client-uuid":   uidStr,             //client uuid
		"x-tt-event-type": "uplink",
	}
	length := Backend_domain.Len()
	if length == 0 {
		if err := OneGetBackdDomain(); err != nil {
			ziLog.Error(fmt.Sprintf("获取后端web服务器失败 %v", err), debug)
			panic("fromGatewayGetForwardMessageHandler 获取后端web服务器失败")
		}
		length = Backend_domain.Len()
		if length == 0 {
			ziLog.Error("后端web为零", debug)
			panic("fromGatewayGetForwardMessageHandler 后端web为零")
		}
	}
	index := int(msgConn.GetId()) % length
	domain := Backend_domain.Get(index)
	if domain == "" {
		if err := OneGetBackdDomain(); err != nil {
			ziLog.Error(fmt.Sprintf("获取后端web服务器失败 %v", err), debug)
			panic("fromGatewayGetForwardMessageHandler 获取后端web服务器失败")
		}
		domain = Backend_domain.Get(0)
		if domain == "" {
			ziLog.Error("获取后端web服务器失败,domain为空", debug)
			panic("fromGatewayGetForwardMessageHandler 获取后端web服务器失败,domain为空")
		}
	}
	backPath, _ := url.JoinPath(domain, extra)
	response, err := common.HttpRespond("POST", backPath, data, header)
	if err != nil {
		ziLog.Error(fmt.Sprintf("fromGatewayGetForwardMessageHandler response err: %v", err), debug)
		// os.Exit(1)
		Backend_domain.Remove(domain)
		if Backend_domain.Len() > 0 {
			time.Sleep(100 * time.Millisecond) // 等待100毫秒后重试
			fromGatewayGetForwardMessageHandler(uidStr, msgConn, data, extra)
		} else {
			if err := OneGetBackdDomain(); err != nil {
				ziLog.Error(fmt.Sprintf("获取后端web服务器失败 %v", err), debug)
				panic("fromGatewayGetForwardMessageHandler 获取后端web服务器失败")
			}
		}
		ziLog.Error(fmt.Sprintf("访问后端web服务器失败, 状态码错误 %v, domain: %v, extra: %v", response.StatusCode, domain, extra), debug)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ziLog.Error(fmt.Sprintf("访问后端web服务器失败, 状态码错误 %v, domain: %v, extra: %v", response.StatusCode, domain, extra), debug)
		panic("fromGatewayGetForwardMessageHandler 访问后端web服务器失败, 状态码错误")
	}

	var (
		result any
	)

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		ziLog.Error(fmt.Sprintf("访问后端web服务器失败 json.NewDecoder %v", err), debug)
		panic(fmt.Sprintf("访问后端web服务器失败 json.NewDecoder %v", err))
	}
	// fmt.Println("node 获取的返回信息", result)
}

// 获取后端web服务器
func OneGetBackdDomain() error {
	respond, err := ectdClient.Client.Get(first_ctx, path.Join("/", config.Project, backend_domain_key), clientv3.WithPrefix(), clientv3.WithPrevKV())
	if err != nil {
		return errors.New("获取web服务器失败: " + err.Error())
	}
	if len(respond.Kvs) == 0 {
		return errors.New("获取web服务器失败: web服务器数量为零")
	}
	for _, kv := range respond.Kvs {
		Backend_domain.Add(string(kv.Value))
		ctx, cancel := context.WithCancel(first_ctx)
		SseCancel.Add(string(kv.Value), cancel)
		backPath, _ := url.JoinPath(string(kv.Value), "sse")
		if len(respond.Kvs) == 1 {
			go SseConn(ctx, backPath)
			go SseConn(ctx, backPath)
		} else {
			go SseConn(ctx, backPath)
		}
		ziLog.Info(fmt.Sprintf("获取web服务器成功 %v %v", string(kv.Key), string(kv.Value)), debug)
	}
	return nil
}
