package register

type Service struct {
	//服务名
	SvrName string          `json:"svr_name"`
	//服务类型
	SvrType int             `json:"svr_type"`
	//节点列表
	SvrNodes       []*Node  `json:"svr_nodes"`
}

//当前节点
type Node struct {
	//机房
	NodeIDC string       `json:"node_idc"`
	//id
	NodeId int           `json:"node_id"`
	//版本
	NodeVersion int      `json:"node_version"`
	//ip地址
	NodeIp string        `json:"node_ip"`
	//端口
	NodePort string      `json:"node_port"`
	//权重
	NodeWeight int       `json:"node_weight"`
	//方法集合
	NodeFuncs []string   `json:"node_funcs"`
}