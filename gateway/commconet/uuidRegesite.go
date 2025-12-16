package commconet

import (
	"context"
	"fmt"
	"path"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/common"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// client gateway uuid regesite
func uuidRegesite(ctx context.Context, msg msginterface.MsgConn, uid uuid.UUID) error {
	if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.Uuid_Register_key, uid.String()), Uuid.String(), clientv3.WithLease(msg.GetLease())); err != nil {
		msg.SetOnline(false)
		msg.Close()
		msg.Cancel()
		ziLog.Error(fmt.Sprintf("发送消息至etcd失败, : %v, err： %v", msg.RemoteAddr().String(), err), debug)
		return err
	}
	return nil
}

// roomId or userId to client uuid
// func UuidRegesiteByRoomIdOrUserId(ctx context.Context, msg msginterface.MsgConn, url, roomIdOrUserId string, uid uuid.UUID) error {
// 	if _, err := EtcdClient.Client.Put(ctx, url+"/"+roomIdOrUserId, Uuid.String()); err != nil {
// 		msg.SetOnline(false)
// 		msg.Close()
// 		msg.Cancel()
// 		log.Println("发送消息至etcd失败, :", msg.RemoteAddr().String(), ", err： ", err)
// 		return err
// 	}
// 	return nil
// }
