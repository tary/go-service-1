package iserver

import (
	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/framework/idata"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"
	"github.com/GA-TECH-SERVER/zeus/framework/msghandler"
)

// IServiceCtrl 上层服务需要实现的接口
type IServiceCtrl interface {
	OnInit() error
	OnTick()
	OnDestroy()
}

// IServiceBase 底层提供的服务基础接口
type IServiceBase interface {
	IEntityPrototype
	msghandler.IRPCHandlers
	IEntities

	//初始化BaseService
	InitBaseService(serviceName string, serviceType idata.ServiceType, ilocal ILocalService) error
	//获取服务ID
	GetSID() uint64
	//获取服务类型
	GetSType() idata.ServiceType
	//获取服务名
	GetSName() string
	//获取服务信息
	GetServiceInfo() *idata.ServiceInfo
	// SetMetadata 设置元数据
	SetMetadata(key, value string)
	// GetMetadata 获取元数据
	GetMetadata(key string) string
}

// IBaseCtrlService 基础服务和上层服务的合集
type IBaseCtrlService interface {
	IServiceCtrl
	IServiceBase
}

// ILocalService 本地服务用到的接口
type ILocalService interface {
	IBaseCtrlService
	// PostCallMsg 把消息投递给服务
	PostCallMsg(msg *msgdef.CallMsg) error
	// PostCallMsgAndWait 把消息投递给服务并等待执行结果
	PostCallMsgAndWait(msg *msgdef.CallMsg) *idata.RetData
}

// IServiceProxy 服务代理接口
type IServiceProxy interface {
	//GetSID 获取服务ID
	GetSID() uint64
	// GetSType 获取服务类型
	GetSType() idata.ServiceType
	// GetAppID 获取服务所属APP ID
	GetAppID() uint64
	// GetServiceInfo 获取ServiceInfo
	GetServiceInfo() *idata.ServiceInfo
	//GetMetadata 获取metadata, 传入key返回value
	GetMetadata(key string) string
	// GetSess 获取服务代理sess
	GetSess() inet.ISession
	// SyncCall 同步调用，等待返回
	SyncCall(retPtr interface{}, methodName string, args ...interface{}) error
	// AsyncCall 异步调用，立即返回
	AsyncCall(methodName string, args ...interface{}) error
	// SendMsg 发送消息给自己的服务器
	SendMsg(msg inet.IMsg) error
	// IsLocal 是否为本地服务
	IsLocal() bool
	// IsValid proxy是否有效
	IsValid() bool
}

// IServiceProxyMgr 服务代理管理器接口
type IServiceProxyMgr interface {
	// GetServiceByID 通过ID获取proxy，如果不存在则返回一个非空但无效的proxy，可以通过proxy的IsValid判断
	GetServiceByID(id uint64) IServiceProxy
	//删除指定id的服务
	DelServiceByID(id uint64)
	// GetRandService 获取特定类型的任意proxy，如果不存在则返回一个非空但无效的proxy，可以通过proxy的IsValid判断
	GetRandService(stype idata.ServiceType) IServiceProxy
	// GetServiceByType 获取指定类型的所有服务
	GetServiceByType(stype idata.ServiceType) []IServiceProxy
	// GetServiceByFunc 根据自定义函数获取自己想要的proxy
	// f为自定义过滤函数，f的第一个参数为此类型的所有服务，f的第二个参数为GetServiceByFunc的第三个参数，用于额外判断
	GetServiceByFunc(stype idata.ServiceType, f func([]IServiceProxy, interface{}) IServiceProxy, data interface{}) IServiceProxy
	// AsyncCallAll 广播异步调用，立即返回
	AsyncCallAll(stype idata.ServiceType, methodName string, args ...interface{}) error
}

var serviceProxyMgrInst IServiceProxyMgr

// GetServiceProxyMgr 获取ServiceProxy管理器
func GetServiceProxyMgr() IServiceProxyMgr {
	return serviceProxyMgrInst
}

// SetServiceProxyMgr 设置ServiceProxy管理器
func SetServiceProxyMgr(ism IServiceProxyMgr) {
	if ism == nil {
		return
	}

	serviceProxyMgrInst = ism
}

// ILocalServiceMgr 服务代理管理器接口
type ILocalServiceMgr interface {
	GetLocalService(sid uint64) ILocalService
}

// localServiceMgrInst 服务管理器实例
var localServiceMgrInst ILocalServiceMgr

// GetLocalServiceMgr 获取LocalService管理器
func GetLocalServiceMgr() ILocalServiceMgr {
	return localServiceMgrInst
}

// SetLocalServiceMgr 设置LocalService管理器
func SetLocalServiceMgr(ilm ILocalServiceMgr) {
	if ilm == nil {
		return
	}

	localServiceMgrInst = ilm
}
