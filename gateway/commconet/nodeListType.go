package commconet

import (
	"sync"

	common "github.com/kongshui/gnw/common"
	msg "github.com/kongshui/gnw/msg/msginterface"
)

var ()

type (
	NodeList struct {
		Gateway    map[string]msg.MsgConn
		Node       map[string]msg.MsgConn
		Client     map[string]msg.MsgConn
		NodeInt    int64
		GatewayInt int64
		ClientInt  int64
		Lock       *sync.RWMutex
	}
	NodeInfoList struct {
		NodeInfo []common.NodeInfo `json:"node_info"` //节点信息列表
		Lock     *sync.RWMutex
	}
)
