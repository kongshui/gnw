package msg

import (
	"context"
	"sync"
	"time"
)

// messageId 自动删除
func (m *message_id_cancel) AutoDelete() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		<-t.C
		if m.Len() == 0 {
			continue
		}
		for id, timeout := range m.id_timeout_map {
			if timeout <= 0 {
				m.Del(id)
			} else {
				m.id_timeout_map[id]--
			}
		}
	}
}

// 添加messageId到map
func (m *message_id_cancel) Add(id uint64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.id_timeout_map[id] = 0
}

// 从map中删除messageId
func (m *message_id_cancel) Del(id uint64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 函数先取消
	m.id_cancel_map[id]()
	delete(m.id_timeout_map, id)
	delete(m.id_cancel_map, id)
}

// 获取长度
func (m *message_id_cancel) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.id_timeout_map)
}

// 创建新的message_id_cancel
func NewMessageIdCancel() *message_id_cancel {
	return &message_id_cancel{
		id_timeout_map: make(map[uint64]int16),
		id_cancel_map:  make(map[uint64]context.CancelFunc),
		lock:           &sync.RWMutex{},
	}
}
