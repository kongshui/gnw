package msg

import (
	"log"
	"sync"
	"time"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

// 路由自动删除
func (r *router_client_map) AutoDelete() {
	t := time.NewTicker(5 * time.Second)
	for {
		<-t.C
		r.CheckConn()
	}
}

// 检查conn路由
func (r *router_client_map) CheckConn() bool {
	if r.ConnLen() == 0 {
		return true
	}
	for uid, msgConn := range r.conn_map {
		if !msgConn.GetOnline() {
			log.Println("路由自动删除:", uid)
			r.DelUid(uid)
			roomId := r.GetRoomIdByUuid(uid)
			if roomId != "" {
				r.DelRoomId(roomId)
			}
			userId := r.GetUserIdByUuid(uid)
			if userId != "" {
				r.DelUserId(userId)
			}
		}
	}
	return false
}

// 添加路由
func (r *router_client_map) AddUid(id string, c msginterface.MsgConn) bool {
	if r.QueryUid(id) {
		return true
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.conn_map[id] = c
	return true
}

// 查询路由
func (r *router_client_map) QueryUid(id string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.conn_map[id]
	return ok
}

// 获取路由
func (r *router_client_map) GetMsgByUuid(id string) msginterface.MsgConn {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.conn_map[id]
}

// 删除路由
func (r *router_client_map) DelUid(id string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.conn_map, id)
	return true
}

// 获取路由数量
func (r *router_client_map) ConnLen() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.conn_map)
}

// 获取所有路由
func (r *router_client_map) GetConnAll() map[string]msginterface.MsgConn {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.conn_map
}

// 添加roomId
func (r *router_client_map) AddRoomId(roomId string, id string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.room_id_map[roomId] = id
	return true
}

// 删除roomid
func (r *router_client_map) DelRoomId(roomId string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.room_id_map, roomId)
	return true
}

// 查询roomid
func (r *router_client_map) QueryRoomId(roomId string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.room_id_map[roomId]
	return ok
}

// 通过roomid获取uuid
func (r *router_client_map) GetRoomId(roomId string) string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.room_id_map[roomId]
}

// 通过uuid获取roomid
func (r *router_client_map) GetRoomIdByUuid(uuid string) string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for roomId, id := range r.room_id_map {
		if id == uuid {
			return roomId
		}
	}
	return ""
}

// 获取roomid长度
func (r *router_client_map) RoomIdLen() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.room_id_map)
}

// 添加userid
func (r *router_client_map) AddUserId(userId string, id string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.user_id_map[userId] = id
	return true
}

// 删除userid
func (r *router_client_map) DelUserId(userId string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.user_id_map, userId)
	return true
}

// 查询userid
func (r *router_client_map) QueryUserId(userId string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.user_id_map[userId]
	return ok
}

// 通过userid获取uuid
func (r *router_client_map) GetUserId(userId string) string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.user_id_map[userId]
}

// 通过uuid获取userid
func (r *router_client_map) GetUserIdByUuid(uuid string) string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for userId, id := range r.user_id_map {
		if id == uuid {
			return userId
		}
	}
	return ""
}

// 获取userid长度
func (r *router_client_map) UserIdLen() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.user_id_map)
}

// 获取所有的路由
func (r *router_client_map) GetConnAllMap() map[string]msginterface.MsgConn {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.conn_map
}

// new
func NewRouterClientMap() *router_client_map {
	return &router_client_map{conn_map: make(map[string]msginterface.MsgConn), room_id_map: make(map[string]string), user_id_map: make(map[string]string), lock: new(sync.RWMutex)}
}
