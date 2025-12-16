package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/kongshui/gnw/gateway/commconet"
	"github.com/kongshui/gnw/gateway/message"

	gcommon "github.com/kongshui/gnw/common"

	"github.com/kongshui/danmu/common"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 获取节点信息
func getNodeList(ctx context.Context) {
	oneGetNodeList()
	respond := etcdClient.Client.Watch(ctx, path.Join("/", config.Project, common.Node_Register_key), clientv3.WithPrefix())
	for wresp := range respond {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				// 处理新增事件
				var data gcommon.NodeInfo
				if err := json.Unmarshal(ev.Kv.Value, &data); err != nil {
					ziLog.Error(fmt.Sprintf("新增节点信息 json.Unmarshal err:%v", err), debug)
					continue
				}
				nodeInfoList.Add(data)
				ziLog.Info(fmt.Sprintf("新增节点信息 %v %v", string(ev.Kv.Key), string(ev.Kv.Value)), debug)
			case clientv3.EventTypeDelete:
				// 处理删除事件
				ziLog.Info(fmt.Sprintf("删除节点信息 %v %v", string(ev.Kv.Key), string(ev.Kv.Value)), debug)
				newSlice := strings.Split(string(ev.Kv.Key), "/")
				uuid := newSlice[len(newSlice)-1]
				nodeInfoList.DelByUuid(uuid)
				message.NodeList.GetNodeByUUid(uuid).SetOnline(false)
			}
		}
	}
}

func oneGetNodeList() error {
	respond, err := etcdClient.Client.Get(context.Background(), path.Join("/", config.Project, common.Node_Register_key), clientv3.WithPrefix(), clientv3.WithPrevKV())
	if err != nil {
		ziLog.Error(fmt.Sprintf("获取node服务器失败 %v", err), debug)
		return errors.New("获取node服务器失败: " + err.Error())
	}
	if len(respond.Kvs) == 0 {
		return nil
	}
	for _, kv := range respond.Kvs {
		if len(kv.Value) != 0 {
			var data gcommon.NodeInfo
			if err := json.Unmarshal(kv.Value, &data); err != nil {
				ziLog.Error(fmt.Sprintf("获取node服务器失败 json.Unmarshal err:%v", err), debug)
				continue
			}
			nodeInfoList.Add(data)
			ziLog.Info(fmt.Sprintf("获取node服务器成功 %v %v", string(kv.Key), string(kv.Value)), debug)
		}
	}
	return nil
}

// 获取直播间绑定信息
func getRoomList(ctx context.Context) {
	time.Sleep(2 * time.Second)
	respond := etcdClient.Client.Watch(ctx, path.Join("/", config.Project, common.RoomInfo_Register_key), clientv3.WithPrefix())
	for wresp := range respond {
		for _, ev := range wresp.Events {
			if len(ev.Kv.Value) == 0 {
				// 删除注册信息
				etcdClient.Client.Delete(ctx, string(ev.Kv.Key))
				continue
			}
			var (
				data common.RoomRegister
			)
			if err := json.Unmarshal(ev.Kv.Value, &data); err != nil {
				ziLog.Error(fmt.Sprintf("json转换失败， err: %v", err), debug)
				// 删除注册信息
				etcdClient.Client.Delete(ctx, string(ev.Kv.Key))
				continue
			}
			// uid := uuid.MustParse(data.Uuid)
			if !commconet.MessageMap.QueryUid(data.Uuid) {
				// 删除注册信息
				etcdClient.Client.Delete(ctx, string(ev.Kv.Key))
				continue
			}
			msgConn := commconet.MessageMap.GetMsgByUuid(data.Uuid)
			msgConn.SetLevel(data.GradeLevel)
			if data.OpenId != "" {
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.OpenId_Register_Uid_key, data.OpenId), data.Uuid, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("OpenId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.Uid_Register_OpenId_key, data.Uuid), data.OpenId, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("OpenId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
			}
			if data.RoomId != "" {
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.RoomId_Register_Uid_key, data.RoomId), data.Uuid, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("RoomId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.Uid_Register_RoomId_key, data.Uuid), data.RoomId, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("RoomId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
				commconet.MessageMap.AddRoomId(data.RoomId, data.Uuid)
				msgConn.SetRoomId(data.RoomId)
			}
			if data.UserId != "" {
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.UserId_Register_Uid_key, data.UserId), data.Uuid, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("UserId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.Uid_Register_UserId_key, data.Uuid), data.UserId, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("UserId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
				commconet.MessageMap.AddUserId(data.UserId, data.Uuid)
				msgConn.SetUserId(data.UserId)
			}
			if data.RoomId != "" && data.OpenId != "" {
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.RoomId_OpenId_Register_key, data.UserId), data.RoomId, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("RoomId_OpenId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
				if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.RoomId_OpenId_Register_key, data.RoomId), data.UserId, clientv3.WithLease(msgConn.GetLease())); err != nil {
					ziLog.Error(fmt.Sprintf("RoomId_OpenId_Register_key 发送消息至etcd失败, uid:%v, err： %v", data.Uuid, err), debug)
				}
			}

			//删除注册信息
			etcdClient.Client.Delete(ctx, string(ev.Kv.Key))
		}
	}
}

// gateway 初始化
func gatewayInit(ctx context.Context) {
	go getNodeList(ctx)
	// go nodeInfoList.AutoDelete()
	go nodeInfoList.AutoAdd(oneGetNodeList)
	go message.NodeList.CheckClientOnlineLoop(ctx)
	go message.NodeList.CheckNodeOnlineLoop(ctx)
	go getRoomList(ctx)
	go endService(ctx)
}

func endService(ctx context.Context) {
	t := time.NewTicker(3 * time.Second)
	defer message.NodeList.DisconnectAllNode()
	defer t.Stop()
	for {
		select {
		case <-t.C:
			message.NodeList.NodeListSet(ctx, &nodeInfoList.NodeInfo, message.Handler)
			// go message.NodeList.GatewayConnectNode(message.Handler)
		case <-ctx.Done():
			return
		}
	}
}
