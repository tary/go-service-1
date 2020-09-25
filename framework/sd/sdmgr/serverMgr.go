package sdmgr

import (
	"fmt"

	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/net/client"
	"github.com/giant-tech/go-service/framework/net/inet"
	"github.com/giant-tech/go-service/framework/sd/sdsess"
)

// IClient client接口
type IClient = inet.ISession

// IConnect connect接口
type IConnect interface {
	Connect(string, uint64) (*client.Session, error)
}

// svrInfo 服务器信息
type svrInfo struct {
	svrID   uint64
	svrAddr string
	svrSess IClient
}

// ServerMgr server mgr
type ServerMgr struct {
	srvType   int32
	connector IConnect
	stop      bool
}

//Init srvType:需要获取的服务器类型
func (mgr *ServerMgr) Init(srvType int32, connector IConnect) error {
	mgr.srvType = srvType
	mgr.connector = connector

	if connector == nil {
		return fmt.Errorf("connector is nil")
	}

	return nil
}

// GetWatchedServerType 获取服务类型
func (mgr *ServerMgr) GetWatchedServerType() int32 {
	return mgr.srvType
}

// GetRandCli 随机获取一个连接
func (mgr *ServerMgr) GetRandCli() (uint64, inet.ISessionBase, error) {
	return sdsess.GetRandSession(mgr.srvType)
}

// GetCliByID 获取client session
func (mgr *ServerMgr) GetCliByID(srvID uint64) (inet.ISessionBase, error) {
	return sdsess.GetSession(srvID)
}

//BroadcastMsg 广播消息
func (mgr *ServerMgr) BroadcastMsg(msg imsg.IMsg) {
	// mgr.srvIDMap.Range(func(k, v interface{}) bool {
	// 	cli := v.(*svrInfo).svrSess
	// 	cli.Send(msg)
	// 	return true
	// })
}
