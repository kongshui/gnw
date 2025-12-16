package msg

// id_type_map 注册id和type的映射
func (m id_type_map) Register(id uint32, t string) {
	m[id] = t
}

// 通过id获取type
func (m id_type_map) GetType(id uint32) string {
	return m[id]
}

// 通过type获取id
func (m id_type_map) GetId(t string) uint32 {
	for k, v := range m {
		if v == t {
			return k
		}
	}
	return 0
}

// NewRouter_client_map 创建router_client_map
func NewIdTypeMap() *id_type_map {
	return &id_type_map{}
}
