package msg

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	TcpConn struct {
		id             int64              //节点id
		uuid           string             //节点uuid
		name           string             //节点名称
		addr           string             //节点地址
		port           string             //节点端口
		node_type      int8               //节点类型
		state          int8               //节点状态, 0:离线, 1:在线, 2:连接中, 3:连接失败, 4:连接超时, 5:连接断开, 6:连接异常, 7:连接成功, 8:连接关闭,
		online         bool               //节点是否在线
		conn           net.Conn           //节点连接
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

// var (
// 	msgBodyPool sync.Pool = sync.Pool{New: func() any {
// 		return &pmsg.MessageBody{}
// 	}}
// )

// SetGroupId 设置分组id
func (msg *TcpConn) SetGroupId(groupId string) {
	msg.groupId = groupId
}

// GetGroupId 获取分组id
func (msg *TcpConn) GetGroupId() string {
	return msg.groupId
}

// setUUid 设置节点uuid
func (msg *TcpConn) SetUuid(uuid string) {
	msg.uuid = uuid
}

// GetUuid 获取节点uuid
func (msg *TcpConn) GetUuid() string {
	return msg.uuid
}

// SetCtx 设置ctx
func (msg *TcpConn) SetCtx(ctx context.Context) {
	msg.ctx = ctx
}

// GetCtx 获取ctx
func (msg *TcpConn) GetCtx() context.Context {
	return msg.ctx
}

// Calcel 设置取消函数
func (msg *TcpConn) SetCancel(cancel context.CancelFunc) {
	msg.cancel = cancel
}

// Cancel 获取取消函数
func (msg *TcpConn) Cancel() {
	msg.cancel()
}

// SetConn 设置连接
func (msg *TcpConn) SetConn(conn net.Conn) {
	msg.conn = conn
}

// GetConn 获取连接
func (msg *TcpConn) GetConn() net.Conn {
	return msg.conn
}

// Close 关闭连接
func (msg *TcpConn) Close() error {
	return msg.conn.Close()
}

// Read 读取数据
func (msg *TcpConn) Read(b []byte) (n int, err error) {
	return msg.conn.Read(b)
}

// Write 写入数据
func (msg *TcpConn) Write(b []byte) (n int, err error) {
	return msg.conn.Write(b)
}

// LocalAddr 获取本地地址
func (msg *TcpConn) LocalAddr() net.Addr {
	return msg.conn.LocalAddr()
}

// RemoteAddr 获取远程地址
func (msg *TcpConn) RemoteAddr() net.Addr {
	return msg.conn.RemoteAddr()
}

// SetDeadline 设置超时时间
func (msg *TcpConn) SetDeadline(t time.Time) error {
	return msg.conn.SetDeadline(t)
}

// SetReadDeadline 设置读取超时时间
func (msg *TcpConn) SetReadDeadline(t time.Time) error {
	return msg.conn.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时时间
func (msg *TcpConn) SetWriteDeadline(t time.Time) error {
	return msg.conn.SetWriteDeadline(t)
}

// 设置lock
func (msg *TcpConn) SetLock(lock *sync.RWMutex) {
	msg.lock = lock
}

// 获取lock
func (msg *TcpConn) GetLock() *sync.RWMutex {
	return msg.lock
}

// 设置节点类型
func (msg *TcpConn) SetNodeType(t int8) {
	msg.node_type = t
}

// 获取节点类型
func (msg *TcpConn) GetNodeType() int8 {
	return msg.node_type
}

// 设置节点状态, 0:未初始化, 1:在线, 2:连接中, 3:连接失败, 4:连接超时, 5:连接断开, 6:连接异常, 7:连接成功, 8:连接关闭,
func (msg *TcpConn) SetState(s int8) {
	msg.state = s
}

// 获取节点状态, 0:未初始化, 1:在线, 2:连接中, 3:连接失败, 4:连接超时, 5:连接断开, 6:连接异常, 7:连接成功, 8:连接关闭,
func (msg *TcpConn) GetState() int8 {
	return msg.state
}

// 设置节点是否在线
func (msg *TcpConn) SetOnline(o bool) {
	msg.online = o
}

// 获取节点是否在线
func (msg *TcpConn) GetOnline() bool {
	return msg.online
}

// 设置节点名称
func (msg *TcpConn) SetName(name string) {
	msg.name = name
}

// 获取节点名称
func (msg *TcpConn) GetName() string {
	return msg.name
}

// 设置节点地址
func (msg *TcpConn) SetAddr(addr string) {
	msg.addr = addr
}

// 获取节点地址
func (msg *TcpConn) GetAddr() string {
	return msg.addr
}

// 设置节点端口
func (msg *TcpConn) SetPort(port string) {
	msg.port = port
}

// 获取节点端口
func (msg *TcpConn) GetPort() string {
	return msg.port
}

// 在处理消息计数器+1
func (msg *TcpConn) CounterAdd() {
	atomic.AddInt64(&msg.counter, 1)
}

// 在处理消息计数器-1
func (msg *TcpConn) CounterSub() {
	atomic.AddInt64(&msg.counter, -1)
}

// 获取处理消息计数器
func (msg *TcpConn) GetCounter() int64 {
	return atomic.LoadInt64(&msg.counter)
}

// 设置负载
func (msg *TcpConn) SetLoad(l float64) {
	msg.load = l
}

// 获取负载
func (msg *TcpConn) GetLoad() float64 {
	return msg.load
}

// 设置节点id
func (msg *TcpConn) SetId(id int64) {
	msg.id = id
}

// 获取节点id
func (msg *TcpConn) GetId() int64 {
	return msg.id
}

// 设置用户id
func (msg *TcpConn) SetUserId(userId string) {
	msg.userId = userId
}

// 获取用户id
func (msg *TcpConn) GetUserId() string {
	return msg.userId
}

// 设置netType
func (msg *TcpConn) SetNetType(netType string) {
	msg.netType = netType
}

// 获取netType
func (msg *TcpConn) GetNetType() string {
	return msg.netType
}

// 设置房间id
func (msg *TcpConn) SetRoomId(roomid string) {
	msg.roomid = roomid
}

// 获取房间id
func (msg *TcpConn) GetRoomId() string {
	return msg.roomid
}

// 设置租约
func (msg *TcpConn) SetLease(lease clientv3.LeaseID) {
	msg.lease = lease
}

// 获取租约
func (msg *TcpConn) GetLease() clientv3.LeaseID {
	return msg.lease
}

// 设置段位
func (msg *TcpConn) SetLevel(level int64) {
	msg.gradeLevel = level
}

// 获取段位
func (msg *TcpConn) GetLevel() int64 {
	return msg.gradeLevel
}

// 获取节点端口
func (msg *TcpConn) Connect() error {
	if msg.netType == "" {
		msg.netType = "tcp"
	}
	if msg.port == "" || msg.addr == "" {
		return errors.New("addr or port is empty")
	}
	c, err := net.Dial(msg.netType, msg.addr+":"+msg.port)
	if err != nil {
		return err
	}
	msg.SetConn(c)
	msg.SetState(7)
	msg.SetOnline(true)
	return nil
}

// NewMsgConn 创建消息连接
func NewMsgConn(conn net.Conn, isHeartBeatOpen bool) *TcpConn {
	ctx, cancel := context.WithCancel(context.Background())
	return &TcpConn{
		conn:          conn,
		lock:          &sync.RWMutex{},
		ctx:           ctx,
		cancel:        cancel,
		netType:       "tcp",
		heartBeatOpen: isHeartBeatOpen,
	}
}

// <KEY> 创建消息连接with ctx和cancel
func NewMsgConnWithCtxAndCancel(ctx context.Context, cancel context.CancelFunc, isHeartBeatOpen bool) *TcpConn {
	return &TcpConn{
		ctx:           ctx,
		cancel:        cancel,
		lock:          &sync.RWMutex{},
		netType:       "tcp",
		heartBeatOpen: isHeartBeatOpen,
	}
}

// NewMsgConnWithCtx 创建消息连接with ctx 和 conn
func NewMsgConnWithCtx(ctx context.Context, cancel context.CancelFunc, conn net.Conn, isHeartBeatOpen bool) *TcpConn {
	return &TcpConn{
		conn:          conn,
		cancel:        cancel,
		ctx:           ctx,
		lock:          &sync.RWMutex{},
		netType:       "tcp",
		heartBeatOpen: isHeartBeatOpen,
	}
}
