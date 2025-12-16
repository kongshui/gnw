package common

type (
	NodeInfo struct {
		Uuid     string `json:"uuid"`      //节点uuid
		Name     string `json:"name"`      //节点名称
		Addr     string `json:"addr"`      //节点地址
		Port     string `json:"port"`      //节点端口
		GroupId  string `json:"group_id"`  //分组id
		NodeType int8   `json:"node_type"` //节点类型
		TTL      int    `json:"ttl"`       //节点过期时间
	}
)
