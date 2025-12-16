package msg

import (
	"context"
	"net"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	MsgConn interface {
		Close() error
		Read(b []byte) (n int, err error)
		Write(b []byte) (n int, err error)
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
		SetGroupId(gid string)
		GetGroupId() string
		SetCtx(ctx context.Context)
		GetCtx() context.Context
		SetCancel(cancel context.CancelFunc)
		Cancel()
		SetName(name string)
		GetName() string
		SetDeadline(t time.Time) error
		SetReadDeadline(t time.Time) error
		SetWriteDeadline(t time.Time) error
		SetLock(lock *sync.RWMutex)
		GetLock() *sync.RWMutex
		GetCounter() int64
		SetNodeType(t int8)
		GetNodeType() int8
		SetState(s int8)
		GetState() int8
		SetOnline(o bool)
		GetOnline() bool
		SetAddr(addr string)
		GetAddr() string
		SetPort(port string)
		GetPort() string
		CounterAdd()
		CounterSub()
		SetLoad(l float64)
		GetLoad() float64
		SetId(id int64)
		GetId() int64
		SetUuid(uuid string)
		GetUuid() string
		SetUserId(userId string)
		GetUserId() string
		ReceiveMessage(handler map[uint32]func(string, MsgConn, []byte, string))
		MessageWrite(msg []byte) (n int, err error)
		Connect() error
		SetNetType(netType string)
		GetNetType() string
		SetRoomId(roomid string)
		GetRoomId() string
		SetLease(lease clientv3.LeaseID)
		GetLease() clientv3.LeaseID
		SetLevel(level int64)
		GetLevel() int64
	}
)
