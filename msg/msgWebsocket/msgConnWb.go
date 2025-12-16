package msg

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	WsConn struct {
		id             int64              //节点id
		uuid           string             //节点uuid
		name           string             //节点名称
		addr           string             //节点地址
		port           string             //节点端口
		node_type      int8               //节点类型
		state          int8               //节点状态, 0:离线, 1:在线, 2:连接中, 3:连接失败, 4:连接超时, 5:连接断开, 6:连接异常, 7:连接成功, 8:连接关闭,
		online         bool               //节点是否在线
		conn           *websocket.Conn    //节点连接
		groupId        string             //分组id
		ctx            context.Context    //ctx
		cancel         context.CancelFunc //取消函数
		lock           *sync.RWMutex      //锁
		netType        string             //网络类型
		counter        int64              //计数器
		load           float64            //负载
		userId         string             //用户id
		roomid         string             //房间id
		lease          clientv3.LeaseID   //租约
		gradeLevel     int64              //段位
		heartBeatCount int32              //心跳次数 等于未接收到心跳次数，时间长了则断开
		heartBeatOpen  bool               //是否保持心跳
	}
)

// SetGroupId 设置分组id
func (msg *WsConn) SetGroupId(groupId string) {
	msg.groupId = groupId
}

// setUUid 设置节点uuid
func (msg *WsConn) SetUuid(uuid string) {
	msg.uuid = uuid
}

// GetUuid 获取节点uuid
func (msg *WsConn) GetUuid() string {
	return msg.uuid
}

// GetGroupId 获取分组id
func (msg *WsConn) GetGroupId() string {
	return msg.groupId
}

// SetCtx 设置ctx
func (msg *WsConn) SetCtx(ctx context.Context) {
	msg.ctx = ctx
}

// GetCtx 获取ctx
func (msg *WsConn) GetCtx() context.Context {
	return msg.ctx
}

// Calcel 设置取消函数
func (msg *WsConn) SetCancel(cancel context.CancelFunc) {
	msg.cancel = cancel
}

// Cancel 获取取消函数
func (msg *WsConn) Cancel() {
	msg.cancel()
}

// SetConn 设置连接
func (msg *WsConn) SetConn(conn *websocket.Conn) {
	msg.conn = conn
}

// GetConn 获取连接
func (msg *WsConn) GetConn() *websocket.Conn {
	return msg.conn
}

// Close 关闭连接
func (msg *WsConn) Close() error {
	return msg.conn.Close()
}

// Read 读取数据
func (msg *WsConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

// Write 写入数据
func (msg *WsConn) Write(b []byte) (_ int, err error) {
	msg.GetLock().Lock()
	err = msg.conn.WriteMessage(websocket.BinaryMessage, b)
	msg.GetLock().Unlock()
	return
}

// LocalAddr 获取本地地址
func (msg *WsConn) LocalAddr() net.Addr {
	return msg.conn.LocalAddr()
}

// RemoteAddr 获取远程地址
func (msg *WsConn) RemoteAddr() net.Addr {
	return msg.conn.RemoteAddr()
}

// SetDeadline 设置超时时间
func (msg *WsConn) SetDeadline(t time.Time) error {
	return msg.conn.NetConn().SetDeadline(t)
}

// SetReadDeadline 设置读取超时时间
func (msg *WsConn) SetReadDeadline(t time.Time) error {
	return msg.conn.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时时间
func (msg *WsConn) SetWriteDeadline(t time.Time) error {
	return msg.conn.SetWriteDeadline(t)
}

// 设置lock
func (msg *WsConn) SetLock(lock *sync.RWMutex) {
	msg.lock = lock
}

// 获取lock
func (msg *WsConn) GetLock() *sync.RWMutex {
	return msg.lock
}

// 设置节点类型
func (msg *WsConn) SetNodeType(t int8) {
	msg.node_type = t
}

// 获取节点类型
func (msg *WsConn) GetNodeType() int8 {
	return msg.node_type
}

// 设置节点状态, 0:未初始化, 1:在线, 2:连接中, 3:连接失败, 4:连接超时, 5:连接断开, 6:连接异常, 7:连接成功, 8:连接关闭,
func (msg *WsConn) SetState(s int8) {
	msg.state = s
}

// 获取节点状态, 0:未初始化, 1:在线, 2:连接中, 3:连接失败, 4:连接超时, 5:连接断开, 6:连接异常, 7:连接成功, 8:连接关闭,
func (msg *WsConn) GetState() int8 {
	return msg.state
}

// 设置节点是否在线
func (msg *WsConn) SetOnline(o bool) {
	msg.online = o
}

// 获取节点是否在线
func (msg *WsConn) GetOnline() bool {
	return msg.online
}

// 设置节点名称
func (msg *WsConn) SetName(name string) {
	msg.name = name
}

// 获取节点名称
func (msg *WsConn) GetName() string {
	return msg.name
}

// 设置节点地址
func (msg *WsConn) SetAddr(addr string) {
	msg.addr = addr
}

// 获取节点地址
func (msg *WsConn) GetAddr() string {
	return msg.addr
}

// 设置节点端口
func (msg *WsConn) SetPort(port string) {
	msg.port = port
}

// 获取节点端口
func (msg *WsConn) GetPort() string {
	return msg.port
}

// 在处理消息计数器+1
func (msg *WsConn) CounterAdd() {
	atomic.AddInt64(&msg.counter, 1)
}

// 在处理消息计数器-1
func (msg *WsConn) CounterSub() {
	atomic.AddInt64(&msg.counter, -1)
}

// 获取处理消息计数器
func (msg *WsConn) GetCounter() int64 {
	return atomic.LoadInt64(&msg.counter)
}

// 设置负载
func (msg *WsConn) SetLoad(l float64) {
	msg.load = l
}

// 获取负载
func (msg *WsConn) GetLoad() float64 {
	return msg.load
}

// 设置节点id
func (msg *WsConn) SetId(id int64) {
	msg.id = id
}

// 获取节点id
func (msg *WsConn) GetId() int64 {
	return msg.id
}

// 设置用户id
func (msg *WsConn) SetUserId(userId string) {
	msg.userId = userId
}

// 获取用户id
func (msg *WsConn) GetUserId() string {
	return msg.userId
}

// 设置netType
func (msg *WsConn) SetNetType(netType string) {
	msg.netType = netType
}

// 获取netType
func (msg *WsConn) GetNetType() string {
	return msg.netType
}

// 设置房间id
func (msg *WsConn) SetRoomId(roomid string) {
	msg.roomid = roomid
}

// 获取房间id
func (msg *WsConn) GetRoomId() string {
	return msg.roomid
}

// 设置租约
func (msg *WsConn) SetLease(lease clientv3.LeaseID) {
	msg.lease = lease
}

// 获取租约
func (msg *WsConn) GetLease() clientv3.LeaseID {
	return msg.lease
}

// 设置段位
func (msg *WsConn) SetLevel(level int64) {
	msg.gradeLevel = level
}

// 获取段位
func (msg *WsConn) GetLevel() int64 {
	return msg.gradeLevel
}

// NewMsgConn 创建消息连接
func NewMsgConn(conn *websocket.Conn, isHeartBeatOpen bool) *WsConn {
	ctx, cancel := context.WithCancel(context.Background())
	return &WsConn{
		conn:          conn,
		lock:          &sync.RWMutex{},
		ctx:           ctx,
		cancel:        cancel,
		netType:       "websocket",
		heartBeatOpen: isHeartBeatOpen,
	}
}

// <KEY> 创建消息连接with ctx和cancel
func NewMsgConnWithCtxAndCancel(ctx context.Context, cancel context.CancelFunc, isHeartBeatOpen bool) *WsConn {
	return &WsConn{
		ctx:           ctx,
		cancel:        cancel,
		lock:          &sync.RWMutex{},
		netType:       "websocket",
		heartBeatOpen: isHeartBeatOpen,
	}
}

// NewMsgConnWithCtx 创建消息连接with ctx 和 conn
func NewMsgConnWithCtx(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, isHeartBeatOpen bool) *WsConn {
	return &WsConn{
		conn:          conn,
		cancel:        cancel,
		ctx:           ctx,
		lock:          &sync.RWMutex{},
		netType:       "websocket",
		heartBeatOpen: isHeartBeatOpen,
	}
}

// 链接
func (c *WsConn) Connect() error {
	return nil
}
