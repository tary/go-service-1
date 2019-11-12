package space

import (
	"zeus/sess"
)

// SpaceSesses 场景服务器专用的Sess 管理器
type SpaceSesses struct {
	clientSrv sess.IMsgServer
}

// NewSpaceSesses 创建一个新的Space的Sess管理器
func NewSpaceSesses(protocal string, addr string, maxConns int) *SpaceSesses {
	return &SpaceSesses{
		clientSrv: sess.NewMsgServer(protocal, addr, maxConns),
	}
}

// Init 初始化
func (srv *SpaceSesses) Init() error {

	if err := srv.clientSrv.Start(); err != nil {
		panic(err)
	}

	srv.clientSrv.RegMsgProc(&SpaceSessesMsgProc{srv: srv})

	return nil
}

func (srv *SpaceSesses) MainLoop() {
	srv.clientSrv.MainLoop()
}

// Destroy 退出时调用
func (srv *SpaceSesses) Destroy() {
	srv.clientSrv.Close()
}

// SetEncryptEnabled
func (srv *SpaceSesses) SetEncryptEnabled() {
	srv.clientSrv.SetEncryptEnabled()
}
