package servicediscovery

//ServerInfo 服务器信息
type ServerInfo struct {
	ServerID     uint64 `redis:"serverid"`  //服务器ID
	Type         uint8  `redis:"type"`      //服务器类型
	SrvOuterAddr string `redis:"outeraddr"` //服务器外网地址
	Load         int    `redis:"load"`      //服务器当前负载
}

// ServerList 服务器列表
type ServerList []*ServerInfo

// Len 获取服务器列表长度
func (list ServerList) Len() int {
	return len(list)
}

// Swap 交换
func (list ServerList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// Less 是否小
func (list ServerList) Less(i, j int) bool {
	return list[i].Load < list[j].Load
}

// GetServerID 拿到Server id
func (u *ServerInfo) GetServerID() uint64 {
	return u.ServerID
}

// GetServerType 拿到server type
func (u *ServerInfo) GetServerType() uint8 {
	return u.Type
}

// GetServerAddr 获得server addr
func (u *ServerInfo) GetServerAddr() string {
	return u.SrvOuterAddr
}

// GetLoad 获得server load
func (u *ServerInfo) GetLoad() int {
	return u.Load
}
