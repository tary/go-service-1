package service

import (
	"fmt"

	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/iserver"

	"github.com/cihub/seelog"

	"github.com/giant-tech/go-service/base/serializer"
	"github.com/giant-tech/go-service/framework/msgdef"
)

// LenStackBuf LenStackBuf
const LenStackBuf int = 10240

// localServiceMgr localServiceMgr
var localServiceMgr *LocalServiceMgr

// init init
func init() {
	localServiceMgr = &LocalServiceMgr{
		idServices:   make(map[uint64]*LocalService),
		typeServices: make(map[idata.ServiceType]*LocalService),
	}

	iserver.SetLocalServiceMgr(localServiceMgr)
}

// GetLocalServiceMgr 获取本地服务器管理器
func GetLocalServiceMgr() *LocalServiceMgr {
	return localServiceMgr
}

// LocalServiceMgr 本地服务管理器
type LocalServiceMgr struct {
	localServices []*LocalService
	idServices    map[uint64]*LocalService            // serviceID
	typeServices  map[idata.ServiceType]*LocalService // serviceType
}

// InitLocalService 初始化本进程服务
func (sm *LocalServiceMgr) InitLocalService(sname string) error {
	s, err := createLocalService(sname)
	if err != nil {
		return err
	}

	if _, ok := sm.typeServices[s.GetSType()]; ok {
		return fmt.Errorf("service already init, name: %s", sname)
	}

	seelog.Debug("App InitLocalService, name: ", sname, " sid: ", s.GetSID())

	sm.localServices = append(sm.localServices, s)
	sm.idServices[s.GetSID()] = s
	sm.typeServices[s.GetSType()] = s

	//把本地服务添加到代理服务管理中
	proxy := CreateServiceProxy(s.GetServiceInfo(), true)
	GetServiceProxyMgr().AddServiceProxy(proxy)

	return nil
}

// GetLocalService 获取本地服务
func (sm *LocalServiceMgr) GetLocalService(sid uint64) iserver.ILocalService {
	s, ok := sm.idServices[sid]
	if ok {
		return s
	}

	return nil
}

// GetLocalServiceByType 根据服务类型获取本地服务
func (sm *LocalServiceMgr) GetLocalServiceByType(stype idata.ServiceType) *LocalService {
	s, ok := sm.typeServices[stype]
	if ok {
		return s
	}

	return nil
}

// RunLocalService 运行本地服务
func (sm *LocalServiceMgr) RunLocalService() error {
	for i := 0; i < len(sm.localServices); i++ {
		s := sm.localServices[i]
		s.wg.Add(1)
		go sm.run(s)
	}

	return nil
}

// Destroy 销毁服务
func (sm *LocalServiceMgr) Destroy() {
	for i := len(sm.localServices) - 1; i >= 0; i-- {
		s := sm.localServices[i]
		s.closeSig <- true
		s.wg.Wait()
		sm.destroy(s)
	}
}

func (sm *LocalServiceMgr) run(s *LocalService) {
	s.Run(s.closeSig)
	s.wg.Done()
}

func (sm *LocalServiceMgr) destroy(m *LocalService) {
	m.OnDestroy()
}

// GetAllLocalService 获取所有本地服务, 参数为SID, 如果为0就是不排除自己
func (sm *LocalServiceMgr) GetAllLocalService(sid uint64) []*idata.ServiceInfo {

	var slist []*idata.ServiceInfo
	//seelog.Debug("GetAllLocalService , sid: ", sid, " slist len=", len(sm.localServices), " localServices = ", sm.localServices)
	for _, s := range sm.localServices {
		if s.GetSID() != sid {
			slist = append(slist, s.GetServiceInfo())
		}
	}

	return slist
}

// 判定一个SID是否在本地服务
/*func (sm *LocalServiceMgr) IsInLocalServiceList(sID uint64) bool {
	for _, s := range sm.localServices {
		if s.GetSID() == sID {
			return true
		}
	}
	return false
}*/

// OnClosed app断开了,抛事件给服务处理
func (sm *LocalServiceMgr) OnClosed(infovec []*idata.ServiceInfo) {
	if len(infovec) == 0 {
		return
	}

	data := serializer.SerializeNew(infovec)
	for _, s := range sm.localServices {
		//s.OnDisconnected(infovec)
		msg := &msgdef.CallMsg{
			SType:      uint8(s.GetSType()),
			MethodName: "Disconnected",
			Params:     data,
		}
		s.PostCallMsg(msg)
	}
}

// OnConnected 连接回调
func (sm *LocalServiceMgr) OnConnected(infovec []*idata.ServiceInfo) {

	data := serializer.SerializeNew(infovec)
	for _, s := range sm.localServices {
		msg := &msgdef.CallMsg{
			SType:      uint8(s.GetSType()),
			MethodName: "Connected",
			Params:     data,
		}
		s.PostCallMsg(msg)
	}
}
