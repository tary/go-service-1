package registry

// Service 微服务对应的所有节点信息
type Service struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Nodes   []*Node `json:"nodes"`
}

// Node 节点信息
type Node struct {
	ID       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}
