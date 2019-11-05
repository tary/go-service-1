package servermgr

import (
	"errors"

	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/framework/idata"
	dbservice "github.com/giant-tech/go-service/framework/logicredis"
)

/*
 服务器管理, 目前由redis实现

 redis中的数据结构: 哈希表
 Key: "server:" + ServerID
 ServerID    	服务器ID
 Type   		服务器类型
 OuterAddress,	服务器外网地址
 InnerAddress,	服务器内网地址
 Load,			服务器当前负载
 Token,			Token

*/

// Servermgr 服务器管理类
type Servermgr struct {
	*LoadBalancer
}

var mgr *Servermgr

// Getservermgr 获取服务器管理实例
func Getservermgr() *Servermgr {
	if mgr == nil {
		mgr = &Servermgr{}
		mgr.LoadBalancer = NewLoadBalancer()
	}
	return mgr
}

// Unregister 将服务器信息从redis中删除
func (mgr *Servermgr) Unregister(info *idata.AppInfo) error {
	if err := dbservice.GetServerUtil(info.AppID).Delete(); err != nil {
		return err
	}
	//iserver.GetSrvInst().FireEvent(serverInfoChannel)
	return nil
}

var errParamNil = errors.New("Param is nil")

// RegState 注册服务器信息
func (mgr *Servermgr) RegState(info *idata.AppInfo) error {
	if info == nil {
		return errParamNil
	}

	util := dbservice.GetServerUtil(info.AppID)

	if util.IsExist() {
		log.Error("Server ID is duplicate!!!!! ", info.AppID)
		// panic("server is is dupblicate")
	}

	if err := util.Register(info); err != nil {
		return err
	}
	return nil
}

// RegServiceState 服务信息
func (mgr *Servermgr) RegServiceState(service *idata.ServiceInfo) error {
	if service == nil {
		return errParamNil
	}

	util := dbservice.GetServiceUtil(service.ServiceID)

	if util.IsExist() {
		log.Error("Service ID is duplicate!!!!! ", service.ServiceID)
		// panic("server is is dupblicate")
	}

	if err := util.Register(service); err != nil {
		return err
	}
	return nil
}

// Update 更新服务器信息
func (mgr *Servermgr) Update(info *idata.AppInfo) error {
	if info == nil {
		return errParamNil
	}

	return dbservice.GetServerUtil(info.AppID).Update(info)
}

// UpdateService 更新服务
func (mgr *Servermgr) UpdateService(service *idata.ServiceInfo) error {
	if service == nil {
		return errParamNil
	}

	return dbservice.GetServiceUtil(service.ServiceID).Update(service)
}

// VerifyServer 验证服务器连接有效性
func (mgr *Servermgr) VerifyServer(uid uint64, token string) bool {
	t, err := dbservice.GetServerUtil(uid).GetToken()
	if err != nil {
		log.Error(err)
		return false
	}

	if t != token {
		return false
	}

	return true
}
