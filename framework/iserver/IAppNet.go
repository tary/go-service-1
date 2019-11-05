package iserver

import "github.com/GA-TECH-SERVER/zeus/base/net/inet"

// IAppNet AppNet接口
type IAppNet interface {
	OnServerConnected(srvID uint64)
	InsertSrvSess(srvID uint64, sess inet.ISession)
	Send(appid uint64, msg inet.IMsg) error
}
