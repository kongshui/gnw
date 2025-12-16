package msg

// // 删除
// func (m *message_timeout) AutoDelete(connRouter *router_client_map) {
// 	t := time.NewTicker(1 * time.Second)
// 	defer t.Stop()
// 	for {
// 		<-t.C
// 		for id, timeout := range *&m.conn_timeout {
// 			if timeout == 0 {
// 				if m.Get_no_delete(id) {
// 					m.no_delete[id]--
// 					if m.no_delete[id] == 0 {
// 						m.Del_no_delete(id)
// 					}
// 					continue
// 				}
// 				m.Del(id)
// 			} else {
// 				m.conn_timeout[id]--
// 			}
// 		}
// 	}
// }

// // 添加route消息id
// func (m *message_timeout) Add(id uint64) {
// 	m.lock.Lock()
// 	defer m.lock.Unlock()
// 	m.conn_timeout[id] = MESSAGE_TIMEOUT
// }

// // 查询是否存在消息id
// func (m *message_timeout) Get(id uint64) bool {
// 	_, ok := m.conn_timeout[id]
// 	return ok
// }

// // 删除route消息id
// func (m *message_timeout) Del(id uint64) {
// 	m.lock.Lock()
// 	defer m.lock.Unlock()
// 	delete(m.conn_timeout, id)
// }

// // 加入no_delete
// func (m *message_timeout) Add_no_delete(id uint64) {
// 	m.lock.Lock()
// 	defer m.lock.Unlock()
// 	m.no_delete[id] = NO_DELETE_MESSAGE_TIMEOUT
// }

// // 查询是否存在no_delete
// func (m *message_timeout) Get_no_delete(id uint64) bool {
// 	_, ok := m.no_delete[id]
// 	return ok
// }

// // 从no_delete删除
// func (m *message_timeout) Del_no_delete(id uint64) {
// 	m.lock.Lock()
// 	defer m.lock.Unlock()
// 	delete(m.no_delete, id)
// }

// // new
// func New_message_timeout() *message_timeout {
// 	return &message_timeout{conn_timeout: make(map[uint64]int), lock: new(sync.RWMutex)}
// }
