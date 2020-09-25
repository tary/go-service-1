package iserver

import (
	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/net/inet"
)

// IAppNet AppNet接口
type IAppNet interface {
	OnServerConnected(srvID uint64)
	InsertSrvSess(srvID uint64, sess inet.ISession)
	Send(appid uint64, msg imsg.IMsg) error
}
