package iserver

import "github.com/GA-TECH-SERVER/zeus/framework/idata"

// IApp App接口
type IApp interface {
	GetAppID() uint64
	GetSeq() uint64
	GetAppNet() IAppNet
	GetNotConnectServices() map[idata.ServiceType]idata.ServiceType
	AddPendingCall(*idata.PendingCall)
	DelPendingCall(seq uint64)
	GetPendingCall(seq uint64) *idata.PendingCall
}

var appInst IApp

// GetApp 获取当前app
func GetApp() IApp {
	return appInst
}

// SetApp 设置app单例
func SetApp(app IApp) {
	if app == nil {
		return
	}

	appInst = app
}
