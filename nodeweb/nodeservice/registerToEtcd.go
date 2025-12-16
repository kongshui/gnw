package nodeservice

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	gcommon "github.com/kongshui/gnw/common"

	"github.com/kongshui/danmu/common"

	clientv3 "go.etcd.io/etcd/client/v3"
	// "go.etcd.io/etcd/clientv3"
)

// RegisterToEtcd
func RegisterToEtcd(ctx context.Context) {
	time.Sleep(2 * time.Second)
	var data gcommon.NodeInfo
	data.Uuid = nodeUuid.String()
	setInitData(&data)
	sData, err := json.Marshal(data)
	if err != nil {
		ziLog.Error(fmt.Sprintf("node 初始化etcd失败, err:%v", err), debug)
		os.Exit(1)
	}
	// 前端地址
	// domain := "http://" + config.Server.Addr + ":" + config.Server.WebPort

	// 租约
	listenId := etcdClient.NewLease(ctx, 5)
	if listenId == 0 {
		ziLog.Error("创建租约失败", debug)
		os.Exit(10)
	}
	etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.Node_Register_key, data.Uuid), string(sData), clientv3.WithLease(listenId))
	time.Sleep(3 * time.Second)
	go func(ctx context.Context) {
		//定时发送
		t := time.NewTicker(3 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				_, err = etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.Node_Register_key, data.Uuid), string(sData), clientv3.WithLease(listenId))
				if err != nil {
					ziLog.Error(fmt.Sprintf("发送消息至etcd Node_Register_key 失败, err:%v", err), debug)
				}
				// _, err = etcdClient.Client.Put(ctx, path.Join("/", config.Project, forward_domain_key, config.Server.Addr+":"+config.Server.WebPort), domain, clientv3.WithLease(listenId))
				// if err != nil {
				// 	ziLog.Error(fmt.Sprintf("发送消息至etcd forward_domain_key 失败, err:%v", err), debug)
				// }
				return
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
	etcdClient.KeepLease(ctx, listenId)
	//立即发送
	// _, err = etcdClient.Client.Put(ctx, common.Node_Register_key+"/"+data.Uuid, string(sData), clientv3.WithLease(listenId))
	// if err != nil {
	// 	log.Println("发送消息至etcd失败,3s 后重试， err：", err)
	// }
}

// func getBackDomain(ctx context.Context) {
// 	time.Sleep(1 * time.Second)
// 	message.OneGetBackdDomain()
// 	respond := etcdClient.Client.Watch(ctx, path.Join("/", config.Project, backend_domain_key), clientv3.WithPrefix())
// 	for wresp := range respond {
// 		for _, ev := range wresp.Events {
// 			// fmt.Println(string(ev.Kv.Key), string(ev.Kv.Value), count)
// 			switch ev.Type {
// 			case clientv3.EventTypePut:
// 				// 处理新增事件
// 				ziLog.Info(fmt.Sprintf("新增后端http节点信息, key:%v, value:%v", string(ev.Kv.Key), string(ev.Kv.Value)), debug)
// 				message.Backend_domain.Add(string(ev.Kv.Value))
// 				pctx, cancel := context.WithCancel(ctx)
// 				message.SseCancel.Add(string(ev.Kv.Value), cancel)
// 				backPath, _ := url.JoinPath(string(ev.Kv.Value), "sse")
// 				if message.SseCancel.Len() == 1 {
// 					go message.SseConn(pctx, backPath)
// 					go message.SseConn(pctx, backPath)
// 				} else {
// 					go message.SseConn(pctx, backPath)
// 				}
// 			case clientv3.EventTypeDelete:
// 				// 处理删除事件
// 				key := path.Base(string(ev.Kv.Key))
// 				key = "http://" + key
// 				ziLog.Info(fmt.Sprintf("删除后端http节点信息, key:%v", key), debug)
// 				message.Backend_domain.Remove(key)
// 				message.SseCancel.Delete(key)
// 				// time.Sleep(1 * time.Second)
// 				// message.OneGetBackdDomain()
// 				continue
// 			}
// 		}
// 	}
// }

// setInitData
func setInitData(data *gcommon.NodeInfo) {
	data.Addr = config.Server.Addr
	data.GroupId = config.Server.GroupId
	data.Name = config.Server.Name
	data.NodeType = config.Server.NodeType
	data.Port = config.Server.Port
}
