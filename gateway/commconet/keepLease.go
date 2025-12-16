package commconet

import (
	"github.com/kongshui/danmu/common"

	"context"
	"fmt"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 保持租约
func keepLease(ctx context.Context, uid string, leaseId clientv3.LeaseID) {
	defer func() {
		ziLog.Info(fmt.Sprintf("用户 %v 续约解除, 租约ID: %v", uid, leaseId), debug)
		disConnect(uid)
		// 释放租约
		etcdClient.Revoke(context.Background(), leaseId)
	}()
	// 续约租约
	ziLog.Info(fmt.Sprintf("用户 %v 续约成功, 租约ID: %v", uid, leaseId), debug)
	keepResp, err := etcdClient.KeepAlive(ctx, leaseId)
	if err != nil {
		ziLog.Error(fmt.Sprintf("续约失败 %v", err), debug)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-keepResp:
			if !ok {
				return
			}
		}
	}
}

func disConnect(uid string) {
	var (
		roomInfo common.RoomInfo
	)
	respond, err := etcdClient.Client.Get(context.Background(), path.Join("/", config.Project, common.Uid_Register_RoomId_key, uid))
	if err == nil {
		for _, kv := range respond.Kvs {
			roomInfo.RoomId = string(kv.Value)
		}
	}
	respond, err = etcdClient.Client.Get(context.Background(), path.Join("/", config.Project, common.Uid_Register_OpenId_key, uid))
	if err == nil {
		for _, kv := range respond.Kvs {
			roomInfo.UserId = string(kv.Value)
		}
	}
	// 添加房间信息至断联列表
	DisconnectMap.Add(uid, roomInfo)
	ziLog.Info(fmt.Sprintf("用户 %v 已断开连接, 房间ID: %v, 用户ID: %v", uid, roomInfo.RoomId, roomInfo.UserId), debug)
	// 额外清理用户ID关联
	// etcdClient.Client.Delete(context.Background(), path.Join("/", config.Project, common.Uid_Register_UserId_key, uid))
}
