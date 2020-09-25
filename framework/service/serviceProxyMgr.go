package service

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/framework/net/inet"

	"github.com/cihub/seelog"
)

var serviceMgr *SProxyMgr

func init() {
	serviceMgr = createServiceProxyMgr()
	iserver.SetServiceProxyMgr(serviceMgr)
}

// GetServiceProxyMgr 获取service proxy manager
func GetServiceProxyMgr() *SProxyMgr {
	return serviceMgr
}

// SProxyMgr 服务管理器
type SProxyMgr struct {
	mtx     *sync.Mutex
	typeMap map[idata.ServiceType][]iserver.IServiceProxy
	idMap   map[uint64]iserver.IServiceProxy
}

// createServiceProxyMgr 创建服务管理器
func createServiceProxyMgr() *SProxyMgr {
	return &SProxyMgr{
		mtx:     &sync.Mutex{},
		typeMap: make(map[idata.ServiceType][]iserver.IServiceProxy),
		idMap:   make(map[uint64]iserver.IServiceProxy),
	}
}

// AddAppServiceProxy 加入服务
// 参数是sess和sess对应的服务列表
func (sm *SProxyMgr) AddAppServiceProxy(sess inet.ISession, infoList []*idata.ServiceInfo) {
	for _, info := range infoList {
		proxy := &SProxy{}
		proxy.serviceInfo = info
		proxy.Sess = sess
		proxy.isLocal = false

		sm.AddServiceProxy(proxy)
	}
}

// DelServiceByAppID 根据appID删除属于这个app的所有服务
func (sm *SProxyMgr) DelServiceByAppID(appID uint64) {
	var idList []uint64
	seelog.Debug("SProxyMgr, DelServiceByAppID : ", appID)
	sm.mtx.Lock()
	for _, proxy := range sm.idMap {
		if proxy.GetAppID() == appID {
			idList = append(idList, proxy.GetSID())
		}
	}
	sm.mtx.Unlock()

	for _, id := range idList {
		sm.DelServiceByID(id)
	}
}

// AddServiceProxy 加入服务
func (sm *SProxyMgr) AddServiceProxy(s *SProxy) {
	seelog.Debugf("AddServiceProxy, serviceInfo: %v", s.serviceInfo)

	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if _, ok := sm.idMap[s.serviceInfo.ServiceID]; ok {
		seelog.Error("AddServiceProxy, already exist SID: ", s.serviceInfo.ServiceID)
		return
	}

	sm.idMap[s.serviceInfo.ServiceID] = s
	sm.typeMap[s.serviceInfo.Type] = append(sm.typeMap[s.serviceInfo.Type], s)
}

// GetServiceByID 通过ID获取服务
func (sm *SProxyMgr) GetServiceByID(id uint64) iserver.IServiceProxy {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if s, ok := sm.idMap[id]; ok {
		return s
	}

	return &SNilProxy{}
}

// GetServiceInfoByAppID 通过AppID获取ServiceInfo列表
func (sm *SProxyMgr) GetServiceInfoByAppID(appID uint64) []*idata.ServiceInfo {
	var infoList []*idata.ServiceInfo

	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	for _, proxy := range sm.idMap {
		if proxy.GetAppID() == appID && proxy.GetServiceInfo() != nil {
			infoList = append(infoList, proxy.GetServiceInfo())
		}
	}

	return infoList
}

// DelServiceByID 通过ID删除服务
func (sm *SProxyMgr) DelServiceByID(id uint64) {
	seelog.Debug("DelServiceByID, serviceID: ", id)

	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if s, ok := sm.idMap[id]; ok {
		ss := sm.typeMap[s.(*SProxy).serviceInfo.Type]
		for idx, tempService := range ss {
			if tempService.(*SProxy).serviceInfo.ServiceID == id {
				seelog.Debug("DelServiceByID succeed, serviceID: ", id)
				sm.typeMap[s.(*SProxy).serviceInfo.Type] = append(ss[:idx], ss[idx+1:]...)
				break
			}
		}

		delete(sm.idMap, id)
	}
}

// GetRandService 随机获取服务
func (sm *SProxyMgr) GetRandService(stype idata.ServiceType) iserver.IServiceProxy {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if vec, ok := sm.typeMap[stype]; ok {
		if len(vec) > 0 {
			randIdx := rand.Intn(len(vec))
			return vec[randIdx]
		}
	}

	return &SNilProxy{}
}

// GetServiceByFunc 根据自定义函数获取服务
// stype 服务类型
// f 为自定义函数，f的第一个参数为此类型的所有服务，f的第二个参数为额外传入的数据（GetServiceByFunc的第三个参数）
// data 额外数据，调用f时作为f的第二个参数
func (sm *SProxyMgr) GetServiceByFunc(stype idata.ServiceType, f func([]iserver.IServiceProxy, interface{}) iserver.IServiceProxy, data interface{}) iserver.IServiceProxy {
	if f == nil {
		return &SNilProxy{}
	}

	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if vec, ok := sm.typeMap[stype]; ok {
		if len(vec) > 0 {
			return f(vec, data)
		}
	}

	return &SNilProxy{}
}

// GetServiceByType 通过类型获取服务
func (sm *SProxyMgr) GetServiceByType(stype idata.ServiceType) []iserver.IServiceProxy {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if vec, ok := sm.typeMap[stype]; ok {
		return vec
	}

	var services []iserver.IServiceProxy

	return services
}

// AsyncCallAll 异步调用所有指定类型的服务
func (sm *SProxyMgr) AsyncCallAll(stype idata.ServiceType, methodName string, args ...interface{}) error {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if vec, ok := sm.typeMap[stype]; ok {
		for _, s := range vec {
			s.AsyncCall(methodName, args...)
		}

		return nil
	}

	return fmt.Errorf("no server")
}

// OnServiceClosed 远程服务不可用
func (sm *SProxyMgr) OnServiceClosed(appID uint64) {

	//远程不可用的service 集合
	var infovec []*idata.ServiceInfo

	infovec = sm.GetServiceInfoByAppID(appID)

	//seelog.Debug("SProxyMgr,  sm.idMap=", sm.idMap, " infovec = ", infovec)

	//远程服务不可用后，通知上层断后处理
	GetLocalServiceMgr().OnClosed(infovec)
}

// 远程服务不可用
/*func (sm *SProxyMgr) OnServiceClosed() {

	//远程不可用的service 集合
	var infovec []*idata.ServiceInfo
	var info idata.ServiceInfo
	seelog.Debug("SProxyMgr,  sm.idMap=", sm.idMap)
	for id, proxy := range sm.idMap {
		seelog.Debug("serviceid = ", id, ", proxy = ", proxy, " , proxy sess=", proxy.GetSess())
		if proxy.GetSess() == nil {
			//proxy.OnClosed()

			info.ServiceID = id
			info.Type = proxy.GetSType()
			sm.DelServiceByID(id)

			infovec = append(infovec, &info)
		}
	}
	//远程服务不可用后，通知上层断后处理
	GetLocalServiceMgr().OnClosed(infovec)
}
*/
