package balancer

type Node struct {
	//idc
	NodeIDC		string		`json:"node_idc"`
	//id
	NodeId		int			`json:"node_id"`
	//版本
	NodeVersion	int			`json:"node_version"`
	//权重
	NodeWeight	int			`json:"node_weight"`
}
