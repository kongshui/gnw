package commconet

import (
	"fmt"
	"slices"
	"sync"
	"time"

	common "github.com/kongshui/gnw/common"
)

// NodeInfoList 节点信息自动删除
func (n *NodeInfoList) AutoDelete() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		count := 0
		<-t.C
		if len(n.NodeInfo) == 0 {
			continue
		}
		for i, v := range n.NodeInfo {
			if v.TTL <= 0 {
				ziLog.Info(fmt.Sprintf("因为ttl小于0删除node %v", v), debug)
				n.DelByUuid(v.Uuid)
			} else {
				n.NodeInfo[i].TTL -= 1
			}
			count++
		}
	}
}

// NodeInfoList 节点信息自动添加
func (n *NodeInfoList) AutoAdd(function func() error) {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		<-t.C
		if len(n.NodeInfo) == 0 {
			if err := function(); err != nil {
				ziLog.Error(fmt.Sprintf("添加节点信息失败 %v", err), debug)
				continue
			}
			continue
		}
	}
}

// NodeInfoList 添加节点信息
func (n *NodeInfoList) Add(nodeInfo common.NodeInfo) {
	if n.query(nodeInfo.Uuid) {
		n.UpdateTTL(nodeInfo.Uuid, NODE_TTL)
		return
	}
	n.Lock.Lock()
	defer n.Lock.Unlock()
	nodeInfo.TTL = NODE_TTL
	n.NodeInfo = append(n.NodeInfo, nodeInfo)
	ziLog.Info(fmt.Sprintf("添加节点信息 %+v", nodeInfo), debug)
}

// NodeInfoList 删除节点信息
func (n *NodeInfoList) DelByName(name string) {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	for i, v := range n.NodeInfo {
		if v.Name == name {
			n.NodeInfo = slices.Delete(n.NodeInfo, i, i+1)
			break
		}
	}
}

// NodeInfoList 通过uuid删除节点信息
func (n *NodeInfoList) DelByUuid(uuid string) {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	for i, v := range n.NodeInfo {
		if v.Uuid == uuid {
			// if i == len(n.NodeInfo)-1 {
			// 	n.NodeInfo = n.NodeInfo[:i]
			// 	break
			// }
			n.NodeInfo = slices.Delete(n.NodeInfo, i, i+1)
			break
		}
	}
}

// NodeInfoList 获取节点信息
func (n *NodeInfoList) Get() []common.NodeInfo {
	n.Lock.RLock()
	defer n.Lock.RUnlock()
	return n.NodeInfo
}

// NodeInfoList 更新节点信息
func (n *NodeInfoList) Update(nodeInfo common.NodeInfo) {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	for i, v := range n.NodeInfo {
		if v.Uuid == nodeInfo.Uuid {
			n.NodeInfo[i] = nodeInfo
			n.NodeInfo[i].TTL = NODE_TTL
			break
		}
	}
}

// NodeInfoList 更新ttl信息
func (n *NodeInfoList) UpdateTTL(Uuid string, ttl int) {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	for i, v := range n.NodeInfo {
		if v.Uuid == Uuid {
			n.NodeInfo[i].TTL = ttl
			break
		}
	}
}

// NodeInfoList 查询节点信息
func (n *NodeInfoList) query(Uuid string) bool {
	n.Lock.RLock()
	defer n.Lock.RUnlock()
	for _, v := range n.NodeInfo {
		if v.Uuid == Uuid {
			return true
		}
	}
	return false
}

// NodeInfoList new
func NewNodeInfoList() *NodeInfoList {
	return &NodeInfoList{NodeInfo: make([]common.NodeInfo, 0), Lock: &sync.RWMutex{}}
}
