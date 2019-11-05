package idata

// AppInfo app info
type AppInfo struct {
	AppID        uint64 `redis:"appserverid"` //app ID
	Type         uint8  `redis:"type"`        //服务器类型
	OuterAddress string `redis:"outeraddr"`   //服务器外网地址
	InnerAddress string `redis:"inneraddr"`
	Load         int    `redis:"load"` //服务器当前负载
	Token        string `redis:"token"`
	Status       int    `redis:"status"`
}

// AppList app list
type AppList []*AppInfo

// Len 实现Len方法
func (list AppList) Len() int {
	return len(list)
}
