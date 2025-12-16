package msg

import (
	"sync"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

// new message_name
func NewMessageName() *message_name {
	return &message_name{
		conn_name: make(map[string]msginterface.MsgConn),
		lock:      &sync.RWMutex{},
	}
}

// add元素到message_name
func (m *message_name) Add(name string, c msginterface.MsgConn) bool {
	if m.query(name) {
		return true
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.conn_name[name] = c
	return true
}

// del元素
func (m *message_name) Delete(name string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.conn_name, name)
	return true
}

// 通过名称查询元素是否存在
func (m *message_name) Get(name string) msginterface.MsgConn {
	if m.query(name) {
		return m.conn_name[name]
	}
	return nil
}

// 通过那么查询某个元素是否存在
func (m *message_name) query(name string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.conn_name[name]
	return ok
}
