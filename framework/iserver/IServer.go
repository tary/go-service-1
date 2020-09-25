package iserver

// ICtrl 后代服务器类需要继承的接口
type ICtrl interface {
	OnInit() error
	OnDestroy()
}

// IServer 为rpc提供的接口
type IServer interface {
	GetServerType() uint8
	GetServerID() uint64
	GetServerCtrl() ICtrl

	//转发rpc消息
	//ForwardRpcMsg(msg imsg.IMsg) error
	//GetTickMS() time.Duration

}

var srvInst IServer

// GetSrvInst 获取当前服务器
func GetSrvInst() IServer {
	return srvInst
}

// SetSrvInst 设置单例
func SetSrvInst(srv IServer) {
	if srv == nil {
		return
	}

	srvInst = srv
}
