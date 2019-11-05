package internal

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"

	"github.com/GA-TECH-SERVER/zeus/base/net/client"
	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/net/server"
	"github.com/GA-TECH-SERVER/zeus/framework/idata"
	"github.com/GA-TECH-SERVER/zeus/framework/iserver"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"
	"github.com/GA-TECH-SERVER/zeus/framework/servermgr"
	"github.com/GA-TECH-SERVER/zeus/framework/service"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// AppNet app组网
type AppNet struct {
	*server.Server
	srvType   uint8
	addr      string
	outerAddr string //外网地址
	token     string
	status    int
	load      int

	pendingSesses *sync.Map
	//app作为client的session集合
	clientSesses *sync.Map
	srvSessiones *sync.Map

	appInfo *idata.AppInfo
}

// NewAppNet 新建app组网
func NewAppNet(srvType uint8 /*srvID uint64,*/, addr string, outerAddr string) *AppNet {

	appnet := &AppNet{
		srvType:       srvType,
		addr:          addr,
		outerAddr:     outerAddr,
		token:         "",
		pendingSesses: &sync.Map{},
		clientSesses:  &sync.Map{},
		srvSessiones:  &sync.Map{},
	}

	return appnet
}

// init appnet初始化
func (appnet *AppNet) init() error {

	//注册服务信息
	services := service.GetLocalServiceMgr().GetAllLocalService(0)
	for i := 0; i < len(services); i++ {
		appnet.regServiceInfo(services[i])
	}

	if err := appnet.registerSrvInfo(); err != nil {
		return err
	}

	go appnet.refresh()

	return nil
}

// refresh 刷新
func (appnet *AppNet) refresh() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("refresh panic:", err, ", Stack: ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	appnet.RefreshSrvInfo()

	updateticker := time.NewTicker(10 * time.Second)
	defer updateticker.Stop()

	refreshTicker := time.NewTicker(time.Duration(5) * time.Second)
	defer refreshTicker.Stop()

	for {
		select {
		case <-updateticker.C:
			if err := servermgr.Getservermgr().Update(appnet.appInfo); err != nil {
				log.Error(err)
			}

			//更新service
			services := service.GetLocalServiceMgr().GetAllLocalService(0)
			for i := 0; i < len(services); i++ {
				if err := servermgr.Getservermgr().UpdateService(services[i]); err != nil {
					log.Error(err)
				}
			}

		case <-refreshTicker.C:
			appnet.RefreshSrvInfo()
		}
	}
}

/*func (appnet *AppNet) genToken() string {
	curtime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(curtime, 10))
	io.WriteString(h, strconv.FormatUint(iserver.GetApp().GetAppID(), 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}*/

// genToken 获得token
func (appnet *AppNet) genToken() string {
	return fmt.Sprintf("%d", iserver.GetApp().GetAppID())
}

//registerSrvInfo 注册当前服务器信息到redis，包括生成token等等
func (appnet *AppNet) registerSrvInfo() error {
	appnet.token = appnet.genToken()
	appnet.regSrvInfo()
	return nil
}

// regSrvInfo 往redis数据库注册服务器
func (appnet *AppNet) regSrvInfo() {

	appnet.appInfo = &idata.AppInfo{
		AppID:        iserver.GetApp().GetAppID(),
		Type:         appnet.srvType,
		OuterAddress: appnet.outerAddr,
		InnerAddress: appnet.addr,
		Load:         appnet.load,
		Token:        appnet.token,
		Status:       appnet.status,
	}

	succeed := false

	for i := 0; i < 5; i++ {
		if err := servermgr.Getservermgr().RegState(appnet.appInfo); err != nil {
			log.Error("regist server info fail , try again after seconds ", appnet.appInfo)
			time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)

		} else {
			succeed = true
			break
		}
	}

	if !succeed {
		log.Error("server register to srvnet failed ", appnet.appInfo)
		panic("register server failed")
	}

	return
}

// regServiceInfo 注册service信息
func (appnet *AppNet) regServiceInfo(serviceinfo *idata.ServiceInfo) {

	succeed := false
	for i := 0; i < 5; i++ {
		if err := servermgr.Getservermgr().RegServiceState(serviceinfo); err != nil {
			log.Error("regist service info fail , try again after seconds ", serviceinfo)
			time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)

		} else {
			succeed = true
			break
		}
	}

	if !succeed {
		log.Error("service register to srvnet failed ", appnet.appInfo)
		panic("register service failed")
	}

	return
}

//RefreshSrvInfo 从数据库刷新服务器信息
func (appnet *AppNet) RefreshSrvInfo() {

	remoteSrvList, err := servermgr.Getservermgr().GetServerList()
	if err != nil {
		log.Error("fetch server info failed", err)
		return
	}

	for _, srvInfo := range remoteSrvList {
		appnet.tryConnectToSrv(srvInfo)
	}

}

// GetServiceListFromDB 从数据库刷新服务信息
func (appnet *AppNet) GetServiceListFromDB() []*idata.ServiceInfo {

	serviceList, err := servermgr.Getservermgr().GetServiceList()
	if err != nil {
		log.Error("fetch sevice info failed", err)
		return nil
	}
	return serviceList
}

// checkInConnectionList 检测连接列表
func (appnet *AppNet) checkInConnectionList(srvID uint64) bool {

	_, ok := appnet.pendingSesses.Load(srvID)
	if ok {
		log.Warn("Server contecting, waiting response ", srvID)
		return true
	}

	if is, ok := appnet.clientSesses.Load(srvID); ok {
		s := is.(inet.ISession)
		if s.IsClosed() {
			log.Warnf("Server disconnected, ID:%d \n", srvID)
			appnet.clientSesses.Delete(srvID)
			//删除服务,需要删除redis里的服务吗？(redis里的服务超时自动删除)
			//service.GetServiceProxyMgr().DelServiceByAppID(srvID)
		} else {
			return true
		}
	}

	return false
}

// isNeedConnectedToSrv 是否需要连接到服务器
func (appnet *AppNet) isNeedConnectedToSrv(info *idata.AppInfo) bool {
	return appnet.isClientSess(info) && !appnet.checkInConnectionList(info.AppID) &&
		appnet.isConnectService(info)
}

// isClientSess 是否是客户端session
func (appnet *AppNet) isClientSess(info *idata.AppInfo) bool {
	return iserver.GetApp().GetAppID() < info.AppID /*&& srv.srvType != info.Type*/

}

func (appnet *AppNet) isConnectService(info *idata.AppInfo) bool {
	notMap := iserver.GetApp().GetNotConnectServices()
	if len(notMap) == 0 {
		return true
	}

	filteredServiceList, _ := appnet.getServiceList(info.AppID)
	localServices := service.GetLocalServiceMgr().GetAllLocalService(0)

	for _, sa := range filteredServiceList {
		for _, sb := range localServices {
			if appnet.isConnectType(sa.Type, sb.Type) {
				return true
			}
		}
	}

	return false
}

func (appnet *AppNet) isConnectType(ta idata.ServiceType, tb idata.ServiceType) bool {
	notMap := iserver.GetApp().GetNotConnectServices()
	if notMap[ta] == tb || notMap[tb] == ta {
		log.Debug("needn't connect, ta: ", ta, ", tb: ", tb)
		return false
	}

	return true
}

// tryConnectToSrv 尝试连接到server
func (appnet *AppNet) tryConnectToSrv(info *idata.AppInfo) {

	if !appnet.isNeedConnectedToSrv(info) {
		return
	}

	appnet.pendingSesses.Store(info.AppID, nil)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("tryConnectToSrv panic:", err, ", Stack: ", string(debug.Stack()))
				if viper.GetString("Config.Recover") == "0" {
					panic(err)
				}
			}
		}()

		s, err := client.Dial("tcp", info.InnerAddress)
		if err != nil {
			appnet.pendingSesses.Delete(info.AppID)

			log.Errorf("Connect failed. %d to %d Addr %s, error:%v", iserver.GetApp().GetAppID(), info.AppID, info.InnerAddress, err)
			return
		}

		appnet.pendingSesses.Store(info.AppID, s)

		lp := &ProcApp{}
		lp.RegisterMsgProcFunctions(s)

		s.Send(&msgdef.ClientVerifyReq{
			ServerID: iserver.GetApp().GetAppID(),
			Token:    info.Token,
		})
		s.Start()

		log.Info("SrvNet try connect to ", info.AppID)

	}()
}

// getServiceList 获取服务列表
func (appnet *AppNet) getServiceList(srvID uint64) ([]*idata.ServiceInfo, []*idata.ServiceInfo) {

	//获取service list, 找出srvID对应的session列表
	serviceList := appnet.GetServiceListFromDB()
	var filteredServiceList []*idata.ServiceInfo
	var totalServiceList []*idata.ServiceInfo
	for _, v := range serviceList {
		if v.AppID == srvID {
			filteredServiceList = append(filteredServiceList, v)
			totalServiceList = append(totalServiceList, v)
		}

		/*	if v.AppID == appnet.srvID && !service.GetLocalServiceMgr().IsInLocalServiceList(v.ServiceID) {
				totalServiceList = append(totalServiceList, v)
			}
		*/

	}
	return filteredServiceList, totalServiceList
}

// OnServerConnected 连上特定的server回调
func (appnet *AppNet) OnServerConnected(srvID uint64) {
	log.Info("Connected to Server succeed !!!  ", srvID)

	if _, ok := appnet.clientSesses.Load(srvID); ok {
		log.Error("Session existed, server id:", srvID)
		return
	}

	sess, ok := appnet.pendingSesses.Load(srvID)
	if !ok {
		log.Error("Session is pending, server id", srvID)
		return
	}

	appnet.pendingSesses.Delete(srvID)
	appnet.clientSesses.Store(srvID, sess)

	filteredServiceList, totalServiceList := appnet.getServiceList(srvID)

	service.GetServiceProxyMgr().AddAppServiceProxy(sess.(inet.ISession), filteredServiceList)

	//通知App上所有本地服务, 哪些服务可用(包括本地服务)

	service.GetLocalServiceMgr().OnConnected(totalServiceList)

}

// InsertSrvSess 插入服务器session
func (appnet *AppNet) InsertSrvSess(srvID uint64, sess inet.ISession) {
	log.Info("InsertSrvSess  ", srvID)

	if _, ok := appnet.srvSessiones.Load(srvID); ok {

		//删除原来的srv sess, 可能是断开重新连了
		appnet.deleteSrvSess(srvID)
		//删除服务,需要删除redis里的服务吗？(redis里的服务超时自动删除)
		service.GetServiceProxyMgr().DelServiceByAppID(srvID)
	}

	appnet.srvSessiones.Store(srvID, sess)

	filteredServiceList, totalServiceList := appnet.getServiceList(srvID)
	service.GetServiceProxyMgr().AddAppServiceProxy(sess, filteredServiceList)

	//通知App上所有本地服务, 和远程服务建立连接
	service.GetLocalServiceMgr().OnConnected(totalServiceList)

}

// deleteSrvSess 删除 srvsession
func (appnet *AppNet) deleteSrvSess(srvID uint64) {
	log.Debug("delete srv sess   ", srvID)

	appnet.srvSessiones.Delete(srvID)

}

// Send 发送
func (appnet *AppNet) Send(toAppid uint64, msg inet.IMsg) error {

	if iserver.GetApp().GetAppID() < toAppid {
		isess, ok := appnet.clientSesses.Load(iserver.GetApp().GetAppID())
		if ok {
			isess.(inet.ISession).Send(msg)
			return nil
		}
	} else {
		isess, ok := appnet.srvSessiones.Load(iserver.GetApp().GetAppID())
		if ok {
			isess.(inet.ISession).Send(msg)
			return nil
		}
	}
	return fmt.Errorf("SrvNet server %d  Server not existed, id:%d", iserver.GetApp().GetAppID(), toAppid)
}
