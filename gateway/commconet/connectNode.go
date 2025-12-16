package commconet

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	common "github.com/kongshui/gnw/common"
	msgnew "github.com/kongshui/gnw/msg"
	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

//commconet 节点调用初始化

func CommconetInit(mode string, isGroup string) {
	if mode != "" {
		mode_type = mode
	}

	if isGroup != "" {
		is_group = isGroup
	}
}

// 网关节点初始化,是个循环，多次执行
func (node *NodeList) NodeListSet(ctx context.Context, nodeInfo *[]common.NodeInfo, handlers map[uint32]func(string, msginterface.MsgConn, []byte, string)) {
	// fmt.Println("节点信息", nodeInfo)
	// switch mode_type {
	// case "minload":
	// 	go node.gatewaySortNodeByLoad(ctx, is_group != "")
	// case "mincounter":
	// 	go node.gatewaySortNodeByCounter(ctx, is_group != "")
	// }
	for _, v := range *nodeInfo {
		isNew := true
		if v.NodeType == int8(NODE_TYPE_GATEWAY) {
			for _, v1 := range node.Gateway {
				if v1.GetUuid() == v.Uuid {
					isNew = false
					if v.TTL <= 0 {
						v1.SetOnline(false)
					}
					break
				}
			}
		} else if v.NodeType == int8(NODE_TYPE_NODE) {
			for _, v1 := range node.Node {
				if v1 == nil {
					continue
				}
				// fmt.Println(v1.GetName(), v.Name, v1.GetName() == v.Name)
				if v1.GetUuid() == v.Uuid {
					isNew = false
					if v.TTL <= 0 {
						v1.SetOnline(false)
					}
					break
				}
			}
		}
		// fmt.Println(v, 11111, isNew)
		// 如果是新节点
		if isNew {
			nMsg := msgnew.NewMsgConn(nil, true, config.Server.HeartbeatOpen)
			cCtx, cancel := context.WithCancel(ctx)
			setNodeInit(cCtx, cancel, nMsg, v)
			err := nMsg.Connect()
			if err != nil {
				nMsg.SetState(3)
				ziLog.Error(fmt.Sprintf("节点链接失败 %v 节点名称：%v 节点地址：%v 节点端口：%v", err, nMsg.GetName(), nMsg.GetAddr(), nMsg.GetPort()), debug)
				// node.DelNodeByUUid(v.GetUuid())
				continue
			}
			nMsg.SetState(2)
			go nMsg.ReceiveMessage(handlers)
			if err := gatewayInitMsgConn(nMsg, v); err != nil {
				ziLog.Error(fmt.Sprintf("网关初始化节点失败 %v 节点名称：%v 节点地址：%v 节点端口：%v", err, nMsg.GetName(), nMsg.GetAddr(), nMsg.GetPort()), debug)
				nMsg.SetOnline(false)
				nMsg.SetState(3)
				continue
			}
			if v.NodeType == NODE_TYPE_GATEWAY {
				id := atomic.AddInt64(&node.GatewayInt, 1) //node.GatewayInt++
				nMsg.SetId(id)
				node.Lock.Lock()
				node.Gateway[v.Uuid] = nMsg
				node.Lock.Unlock()
			}
		}
	}
}

// 设置节点初始化
func setNodeInit(ctx context.Context, cancel context.CancelFunc, msg msginterface.MsgConn, node common.NodeInfo) {
	msg.SetCtx(ctx)
	msg.SetCancel(cancel)
	msg.SetUuid(node.Uuid)
	msg.SetName(node.Name)
	msg.SetAddr(node.Addr)
	msg.SetPort(node.Port)
	msg.SetGroupId(node.GroupId)
	msg.SetNodeType(node.NodeType)
	msg.SetOnline(true)

	//不为空时添加
	if node.GroupId != "" {
		if !group_map.Add(node.GroupId, msg) {
			ziLog.Error("group_map 添加节点失败 groupId: "+node.GroupId+" uuid: "+node.Uuid+
				" name: "+node.Name+"addr: "+msg.RemoteAddr().String()+" 节点已存在", debug)
		}
	}
	//不为空时添加
	if node.Name != "" {
		name_map.Add(node.Name, msg)
	}

}

// 初始化msgConn
func gatewayInitMsgConn(msg msginterface.MsgConn, node common.NodeInfo) error {
	// msg.NewMsgConn()
	msg.SetName(node.Name)
	msg.SetAddr(node.Addr)
	msg.SetPort(node.Port)
	msg.SetGroupId(node.GroupId)
	msg.SetNodeType(node.NodeType)
	if node.Addr == "" || node.Port == "" {
		msg.SetOnline(false)
		return errors.New("节点不在线或者地址为空或者端口为空")
	}
	msg.SetOnline(true)
	return nil
}

// 网关链接节点并检查，是个循环，多次执行
// func (node *NodeList) GatewayConnectNode(handlers map[uint32]func(string, msginterface.MsgConn, []byte, string)) {
// 	if len(node.Node) == 0 {
// 		return
// 	}
// 	for _, v := range node.Node {
// 		//查看是否是在线状态
// 		if !v.GetOnline() {
// 			// group_map.DeleteNode(v.GetGroupId(), v)
// 			continue
// 		}
// 		//是否是进行初始化
// 		if v.GetState() != 0 {
// 			continue
// 		}
// 		v.SetState(2)
// 		go v.ReceiveMessage(handlers)
// 	}
// }

// 网关Node节点通过任务数排序
// func (node *NodeList) gatewaySortNodeByCounter(ctx context.Context, isGroup bool) {
// 	t := time.NewTicker(1 * time.Second)
// 	defer t.Stop()
// 	for {
// 		select {
// 		case <-t.C:
// 			if len(node.Node) == 0 {
// 				return
// 			}
// 			if isGroup {
// 				group_map.GroupSortByCounter()
// 				continue
// 			}
// 			node.Lock.Lock()
// 			sort.Slice(node.Node, func(i, j int) bool {
// 				return node.Node[i].GetCounter() < node.Node[j].GetCounter()
// 			})
// 			node.Lock.Unlock()
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }

// 网关node节点通过load排序
// func (node *NodeList) gatewaySortNodeByLoad(ctx context.Context, isGroup bool) {
// 	t := time.NewTicker(1 * time.Second)
// 	defer t.Stop()
// 	for {
// 		select {
// 		case <-t.C:
// 			if len(node.Node) == 0 {
// 				return
// 			}
// 			if isGroup {
// 				group_map.GroupSortByLoad()
// 				continue
// 			}
// 			node.Lock.Lock()
// 			sort.Slice(node.Node, func(i, j int) bool {
// 				return node.Node[i].GetLoad() < node.Node[j].GetLoad()
// 			})
// 			node.Lock.Unlock()
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }

// 断开所有节点链接
func (node *NodeList) DisconnectAllNode() {
	for _, v := range node.Gateway {
		v.Close()
		v.SetState(8)
	}
	for _, v := range node.Node {
		v.Close()
		v.SetState(8)
	}
	for _, v := range node.Client {
		v.Close()
		v.SetState(8)
	}
}

// gateway发送消息到指定节点,通过hash,或者其他排序
func (node *NodeList) SendMsgToNode(group string, name string, data []byte, maxRetry int) error {
	maxRetry++
	if maxRetry > 3 {
		return errors.New("SendMsgToNode 最大重试次数超过3次")
	}
	switch is_group {
	case "group":
		return node.sendMsgToNodeByGroup(group, name, data, maxRetry)
	case "name":
		// msg := msgPool.Get().(msginterface.MsgConn)
		// defer msgPool.Put(msg)
		msg := name_map.Get(name)
		if msg == nil {
			return errors.New("节点名称不存在, name: " + name)
		}
		_, err := msg.MessageWrite(data)
		return err
	default:
		return node.sendMsgToNodeDefault(group, name, data, maxRetry)
	}
}

// defalut sendmessage
func (node *NodeList) sendMsgToNodeDefault(group string, name string, data []byte, maxRetry int) error {
	// msg := msgPool.Get().(msginterface.MsgConn)
	// defer msgPool.Put(msg)
	length := len(node.Node)
	if length == 0 {
		ziLog.Error("SendMsgToNode 节点为空", debug)
		return errors.New("SendMsgToNode 节点为空")
	}
	msg := node.GetNode()
	if !msg.GetOnline() || (msg.GetState() != 2 && msg.GetState() != 7) {
		return errors.New("sendMsgToNodeDefault 所有节点都不在线")
	}
	_, err := msg.MessageWrite(data)
	if err != nil {
		ziLog.Error(fmt.Sprintf("sendMsgToNodeDefault 发送失败 %v", err), debug)
		msg.SetOnline(false)
		return node.SendMsgToNode(group, name, data, maxRetry)
	}
	return nil
}

// 通过group发送
func (node *NodeList) sendMsgToNodeByGroup(group string, name string, data []byte, maxRetry int) error {
	// msg := msgPool.Get().(msginterface.MsgConn)
	// defer msgPool.Put(msg)
	msg := group_map.GetMsgByNameWithHash(group)
	if msg == nil || !msg.GetOnline() || (msg.GetState() != 2 && msg.GetState() != 7) {
		if msg == nil {
			return errors.New("节点组不存在, group: " + group)
		}
	}
	_, err := msg.MessageWrite(data)
	if err != nil {
		ziLog.Error(fmt.Sprintf("sendMsgToNodeByGroup 发送失败 %v", err), debug)
		msg.SetOnline(false)
		return node.SendMsgToNode(group, name, data, maxRetry)
	}
	return nil
}

// 发送消息到最小处理任务数节点h或者最小负载节点
// func (node *NodeList) SendMsgToMinTaskOrMinLoadNode(data []byte) error {
// 	_, err := node.Node[0].MessageWrite(data)
// 	return err
// }

// 获取节点数量
func (node *NodeList) GetNodeCount() int {
	return len(node.Node)
}

// // 获取节点名称
// func (node *NodeList) GetNodeName(i int) string {
// 	return node.Node[i].GetName()
// }

// 添加节点
func (node *NodeList) AddNode(msgConn msginterface.MsgConn) {
	node.Lock.Lock()
	defer node.Lock.Unlock()
	id := atomic.AddInt64(&node.NodeInt, 1) //node.NodeInt++
	msgConn.SetId(id)
	// node.Node = append(node.Node, msgConn)
	node.Node[msgConn.GetUuid()] = msgConn
}

// 注册新的NodeList
func NewNodeList() *NodeList {
	return &NodeList{Gateway: make(map[string]msginterface.MsgConn), Node: make(map[string]msginterface.MsgConn), Client: make(map[string]msginterface.MsgConn), GatewayInt: 0, NodeInt: 0, ClientInt: 0, Lock: &sync.RWMutex{}}
}

// 获取节点
func (node *NodeList) GetNodeByUUid(uuid string) msginterface.MsgConn {
	for _, v := range node.Node {
		if v.GetUuid() == uuid {
			return v
		}
	}
	return msgnew.NewMsgConn(nil, true, config.Server.HeartbeatOpen) // 返回一个新的MsgConn对象，避免nil引用错误
}

// 删除节点
func (node *NodeList) DelNodeByUUid(uuid string) error {
	// for i, v := range node.Node {
	// 	if v.GetUuid() == uuid {
	// 		v.SetOnline(false)
	// 		v.Cancel()
	// 		group_map.Delete(v.GetGroupId())
	// 		name_map.Delete(v.GetName())
	// 		node.Node = slices.Delete(node.Node, i, i+1)
	// 		break
	// 	}
	// }
	// atomic.AddInt64(&node.NodeInt, -1)
	if _, exists := node.Node[uuid]; !exists {
		return nil
	}
	// 删除node.Node中的节点
	delete(node.Node, uuid)
	return nil
}

// 检查Node节点是否在线
func (node *NodeList) checkNodeOnline() {
	if len(node.Node) == 0 {
		return
	}
	// fmt.Println("检查Node节点是否在线", len(node.Node))
	node.Lock.Lock()
	defer node.Lock.Unlock()
	for _, v := range node.Node {
		if v == nil {
			continue
		}
		if !v.GetOnline() {
			fmt.Println("因为不在线删除node", v)
			v.Cancel()
			group_map.Delete(v.GetGroupId())
			name_map.Delete(v.GetName())
			node.DelNodeByUUid(v.GetUuid())
			// atomic.AddInt64(&node.ClientInt, -1)
			continue
		}
	}
}

// 持续检查Node节点是否在线
func (node *NodeList) CheckNodeOnlineLoop(ctx context.Context) {
	t := time.NewTicker(3 * time.Second)
	var check int32 = 0
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if check == 0 {
				atomic.AddInt32(&check, 1)
				node.checkNodeOnline()
			}

		case <-ctx.Done():
			return
		}
		atomic.StoreInt32(&check, 0)
	}
}

// 通过下标获取节点
func (node *NodeList) GetNode() msginterface.MsgConn {
	// node.Lock.RLock()
	// defer node.Lock.RUnlock()
	var msg msginterface.MsgConn
	count := 0
	switch mode_type {
	case "minload":
		for _, v := range node.Node {
			if v.GetOnline() && (v.GetState() == 2 || v.GetState() == 7) && count == 0 {
				count++
				msg = v
			} else if v.GetOnline() && (v.GetState() == 2 || v.GetState() == 7) && int64(v.GetLoad()) < int64(msg.GetLoad()) {
				count++
				if count == 1 {
					msg = v
				}
			} else {
				continue
			}
		}
	case "mincounter":
		for _, v := range node.Node {
			if v.GetOnline() && (v.GetState() == 2 || v.GetState() == 7) && count == 0 {
				count++
				msg = v
			} else if v.GetOnline() && (v.GetState() == 2 || v.GetState() == 7) && v.GetCounter() < msg.GetCounter() {
				count++
				if count == 1 {
					msg = v
				}
			} else {
				continue
			}
		}
	default:
		for _, v := range node.Node {
			if v.GetOnline() && (v.GetState() == 2 || v.GetState() == 7) {
				return v
			}
		}
	}
	if msg == nil {
		return msgnew.NewMsgConn(nil, true, config.Server.HeartbeatOpen) // 返回一个新的MsgConn对象，避免nil引用错误
	}
	return msg
	// if index < 0 || index >= len(node.Node) {
	// 	if len(node.Node) > 0 {
	// 		return node.Node[0]
	// 	}
	// 	return msgnew.NewMsgConn(nil, true, true) // 返回一个新的MsgConn对象，避免nil引用错误
	// }
	// return node.Node[index]
}
