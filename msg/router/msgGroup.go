package msg

import (
	"sync"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

// new
func NewMessageGroup() *message_group {
	return &message_group{
		conn_group: make(map[string]map[string]msginterface.MsgConn),
		lock:       &sync.RWMutex{},
	}
}

// 添加到组
func (g *message_group) Add(group string, c msginterface.MsgConn) bool {
	g.lock.Lock()
	defer g.lock.Unlock()
	if _, ok := g.conn_group[group]; !ok {
		g.conn_group[group] = make(map[string]msginterface.MsgConn)
	}
	if _, ok := g.conn_group[group][c.GetUuid()]; !ok {
		g.conn_group[group][c.GetUuid()] = c
		return true
	}

	return false
}

// 删除组
func (g *message_group) Delete(id string) bool {
	g.lock.Lock()
	defer g.lock.Unlock()
	delete(g.conn_group, id)
	return true
}

// 删除组中的节点
func (g *message_group) DeleteNode(group string, c msginterface.MsgConn) bool {
	g.lock.Lock()
	defer g.lock.Unlock()
	_, ok := g.conn_group[group][c.GetUuid()]
	if !ok {
		return false
	}
	delete(g.conn_group[group], c.GetUuid())
	return true
}

// 获取组
func (g *message_group) Get(id string) map[string]msginterface.MsgConn {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.conn_group[id]
}

// 获取所有组
func (g *message_group) GetAll() map[string]map[string]msginterface.MsgConn {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.conn_group
}

// 获取组数量
func (g *message_group) Len() int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return len(g.conn_group)
}

// 获取指定组数量
func (g *message_group) GetLen(id string) int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return len(g.conn_group[id])
}

// 按照名称和下表获取链接
func (g *message_group) GetMsgByNameWithHash(groupId string) msginterface.MsgConn {
	for _, v := range g.conn_group[groupId] {
		if v.GetOnline() && (v.GetState() == 2 || v.GetState() == 7) {
			return v
		}
	}
	return nil
}
