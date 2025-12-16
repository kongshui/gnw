package msg

import (
	"context"
	"sync"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

type (
	router_client_map struct {
		conn_map    map[string]msginterface.MsgConn // uuid和msg结合
		room_id_map map[string]string               // room和uuid结合
		user_id_map map[string]string               // user和uuid结合
		lock        *sync.RWMutex
	}
	// message_timeout struct {
	// 	conn_timeout map[uint64]int
	// 	// no_delete    map[uint64]int
	// 	lock *sync.RWMutex
	// }
	message_group struct {
		conn_group map[string]map[string]msginterface.MsgConn // group和uuid结合
		lock       *sync.RWMutex
	}
	message_name struct {
		conn_name map[string]msginterface.MsgConn
		lock      *sync.RWMutex
	}
	id_type_map map[uint32]string

	message_id_cancel struct {
		id_cancel_map  map[uint64]context.CancelFunc
		id_timeout_map map[uint64]int16
		lock           *sync.RWMutex
	}
)
