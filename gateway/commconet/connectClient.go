package commconet

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"sync/atomic"
	"time"

	msgnew "github.com/kongshui/gnw/msg"
	msginterface "github.com/kongshui/gnw/msg/msginterface"

	gcommon "github.com/kongshui/gnw/common"

	"github.com/kongshui/danmu/common"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

// client 加入网关
func (node *NodeList) GatewayAddClient(ctx context.Context, cancel context.CancelFunc, conn any, handlers map[uint32]func(string, msginterface.MsgConn, []byte, string)) {
	//设置上下文
	msg := msgnew.NewMsgConn(conn, reflect.TypeOf(conn).String() == "*net.Conn", config.Server.HeartbeatOpen)
	id := atomic.AddInt64(&node.ClientInt, 1)
	setNodeInit(ctx, cancel, msg, gcommon.NodeInfo{NodeType: NODE_TYPE_CLIENT})
	msg.SetId(id)
	msg.SetOnline(true)
	msg.SetState(2)
	// 加锁
	node.Lock.Lock()

	uid := uuid.New()
	msg.SetUuid(uid.String())
	// 添加到Client
	node.Client[uid.String()] = msg
	// 设置租约
	leaseId := etcdClient.NewLease(ctx, 10)
	ziLog.Info(fmt.Sprintf("client 加入网关 %v %v", msg.RemoteAddr().String(), uid.String()), debug)
	msg.SetLease(leaseId)
	// 持续租约
	go keepLease(ctx, uid.String(), leaseId)
	// 添加到Online
	if _, err := etcdClient.Client.Put(ctx, path.Join("/", config.Project, common.Uuid_Online_key, uid.String()), uid.String(), clientv3.WithLease(leaseId)); err != nil {
		ziLog.Error(fmt.Sprintf("Uuid_Online_key 发送消息至etcd失败, uid:%v, err： %v", uid, err), debug)
		// 解锁
		node.Lock.Unlock()
		return
	}
	// 添加到消息列表
	MessageMap.AddUid(uid.String(), msg)
	if err := uuidRegesite(ctx, msg, uid); err != nil {
		ziLog.Error(fmt.Sprintf("uuidRegesite error:%v", err), debug)
		// 解锁
		node.Lock.Unlock()
		return
	}
	data := &pmsg.UidSend{UidStr: uid.String()}
	gdata, _ := proto.Marshal(data)
	// log.Println(data.String())
	sData := &pmsg.MessageBody{MsgId: pmsg.MessageId_UidGet, MessageData: gdata}
	sdataByte, err := proto.Marshal(sData)
	if err != nil {
		ziLog.Error(fmt.Sprintf("proto.Marshal sData error:%v", err), debug)
		// 解锁
		node.Lock.Unlock()
		return
	}
	if _, err := msg.MessageWrite(sdataByte); err != nil {
		ziLog.Error(fmt.Sprintf("MessageWrite sdataByte error:%v", err), debug)
		// 解锁
		node.Lock.Unlock()
		return
	}
	// 解锁
	node.Lock.Unlock()
	msg.ReceiveMessage(handlers)
}

// 查询客户端userId是否在线
func (node *NodeList) CheckUserIdOnline(userId string) bool {
	if len(node.Node) == 0 {
		return false
	}
	for _, v := range node.Client {
		if v.GetUserId() == userId {
			return true
		}
	}
	return false
}

// 检查client节点是否在线
func (node *NodeList) checkClientOnline() {
	if len(node.Client) == 0 {
		return
	}
	for K, v := range node.Client {
		if v == nil {
			continue
		}
		if !v.GetOnline() {
			ziLog.Info(fmt.Sprintf("client 节点 %v 已离线, 地址 %v", v.GetUuid(), v.RemoteAddr().String()), debug)
			v.Close()
			v.Cancel()
			delete(node.Client, K)
			if v.GetGroupId() != "" {
				etcdClient.Client.Delete(context.Background(), path.Join("/", config.Project, common.GroupId_Register_key, v.GetUuid()))

			}
			// atomic.AddInt64(&node.ClientInt, -1)
			continue
		}
	}
}

// 持续检查client节点是否在线
func (node *NodeList) CheckClientOnlineLoop(ctx context.Context) {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			node.checkClientOnline()
		case <-ctx.Done():
			return
		}
	}
}
