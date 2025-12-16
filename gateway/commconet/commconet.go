package commconet

import (
	"github.com/kongshui/gnw/gateway/gwinit"
	msg "github.com/kongshui/gnw/msg/router"

	"github.com/kongshui/danmu/common"

	"github.com/google/uuid"
)

var (
	MessageMap    = msg.NewRouterClientMap()
	etcdClient    = gwinit.Ectd_client
	My_Node_Type  int
	mode_type     = "hash"                //默认hash模式hash模式，minload最小负载，mincounter最小任务，group分组
	group_map     = msg.NewMessageGroup() //分组
	is_group      = ""                    // 通过name
	name_map      = msg.NewMessageName()  //通过那么分组
	Uuid          = uuid.New()            //设置uuid
	DisconnectMap = common.NewStringToRoomInfoMap()
	// msgPool       sync.Pool = sync.Pool{
	// 	New: func() any {
	// 		msg := msgnew.NewMsgConn(nil, true, true)
	// 		return msg
	// 	},
	// }
	ziLog  = &gwinit.Zilog // 日志
	debug  = gwinit.Debug  // 调试
	config = gwinit.Config
)

func init() {
	//后续查看怎么处理
	go MessageMap.AutoDelete()
	My_Node_Type = NODE_TYPE_GATEWAY
}
