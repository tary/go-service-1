package internal

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/giant-tech/go-service/base/net/server"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
	"github.com/giant-tech/go-service/base/zlog"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/framework/msgdef"
	"github.com/giant-tech/go-service/framework/servermgr"
	"github.com/giant-tech/go-service/framework/service"

	"github.com/cihub/seelog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
)

// MyApp app实例
var MyApp *App

// init 初始化
func init() {
	MyApp = &App{}

	iserver.SetApp(interface{}(MyApp).(iserver.IApp))
}

// setConfig 设置配置
func setConfig(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败: " + configPath + ", err: " + err.Error())
	}
}

// getServiceNameList 拿到服务名列表
func getServiceNameList(str string) []string {
	// 去除空格
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)

	services := strings.Split(str, ",")

	return services
}

// getNotConnectServiceMap 获取不会连接的服务对
func getNotConnectServiceMap(str string) map[idata.ServiceType]idata.ServiceType {
	tempMap := make(map[idata.ServiceType]idata.ServiceType)

	if len(str) == 0 {
		return tempMap
	}

	// 去除空格
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)

	serviceCouples := strings.Split(str, ",")

	for _, couple := range serviceCouples {
		kv := strings.Split(couple, ":")
		if len(kv) != 2 {
			panic("NotConnectServiceMap incorrect")
		}

		type0 := service.GetServiceByName(kv[0])
		type1 := service.GetServiceByName(kv[1])

		if type0 == nil || type1 == nil {
			panic("NotConnectServiceMap GetServiceByName not found")
		}

		tempMap[type0.ServiceTypeID] = type1.ServiceTypeID
		tempMap[type1.ServiceTypeID] = type0.ServiceTypeID
	}

	return tempMap
}

// App 服务器基类
type App struct {
	//*server.Server
	*AppNet
	appID       uint64
	startupTime time.Time

	pendingMap sync.Map
	seq        atomic.Uint64

	//不会连接的两个服务
	notConnectServices map[idata.ServiceType]idata.ServiceType

	pendingClose chan bool
}

// init app初始化
func (srv *App) init(names ...string) error {
	msgdef.Init()

	srv.startupTime = time.Now()

	var err error
	srv.appID, err = dbservice.GetIDGenerator().GetGlobalID()
	if err != nil {
		seelog.Error("get server id error: ", err)
		return err
	}

	seelog.Debug("App.init, appID: ", srv.appID)

	//初始化服务
	for i := 0; i < len(names); i++ {
		err = service.GetLocalServiceMgr().InitLocalService(names[i])
		if err != nil {
			seelog.Error("InitLocalService error: ", err)
			return err
		}
	}

	listenAddr := viper.GetString("ServerApp.ListenAddr")
	srv.AppNet = NewAppNet(1, listenAddr, "")
	srv.AppNet.init()

	srv.AppNet.Server, err = server.New("tcp", listenAddr, 0)
	if err != nil {
		return err
	}
	//多设了一遍，可以删掉？	srv.AppNet.Server.SetVerifyMsgID(msgdef.ClientVerifyReqMsgID)

	// 添加MsgProc, 这样新连接创建时会注册处理函数
	srv.AppNet.Server.AddMsgProc(&ProcApp{})

	service.GetLocalServiceMgr().RunLocalService()

	return nil
}

// destroy app销毁
func (srv *App) destroy() {
	seelog.Debug("App.destory")

	srv.pendingClose <- true

	service.GetLocalServiceMgr().Destroy()
	//删除redis里的app信息,
	info := &idata.AppInfo{
		AppID: srv.appID,
	}
	servermgr.Getservermgr().Unregister(info)
	// 删除app上的所有service信息
	service.GetServiceProxyMgr().DelServiceByAppID(srv.appID)
	if srv.Server != nil {
		srv.Server.Close()
	}
}

// Run 逻辑入口
// configFile 配置文件
func (srv *App) Run(configFile string) {
	pflag.Uint("pprof-port", 0, "pprof http port")
	pflag.String("configfile", "../res/config/server.toml", "config file")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if configFile == "" {
		configFile = viper.GetString("configfile")
	}

	setConfig(configFile)

	// 设置Seelog
	zlog.InitDefault()
	defer seelog.Flush()

	//start prof
	startProfServer()

	serviceString := viper.GetString("ServerApp.Services")
	services := getServiceNameList(serviceString)
	if len(services) == 0 {
		panic("no service")
	}

	srv.notConnectServices = getNotConnectServiceMap(viper.GetString("ServerApp.NotConnect"))

	seelog.Info("services", services)
	if err := srv.init(services...); err != nil {
		panic(err)
	}

	if srv.Server != nil {
		go srv.Server.Run()
	}

	srv.pendingClose = make(chan bool, 1)
	go srv.loopCheckPendingCall(srv.pendingClose)

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	//signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGSEGV, syscall.SIGPIPE, syscall.SIGTERM)

	<-c

	srv.destroy()
}

// GetAppID 获得appID
func (srv *App) GetAppID() uint64 {
	return srv.appID
}

// GetSeq 获得请求
func (srv *App) GetSeq() uint64 {
	return srv.seq.Inc()
}

// GetNotConnectServices 获取不会连接的服务对
func (srv *App) GetNotConnectServices() map[idata.ServiceType]idata.ServiceType {
	return srv.notConnectServices
}

// AddPendingCall 添加等待调用
func (srv *App) AddPendingCall(call *idata.PendingCall) {
	//seelog.Debug("AddPendingCall, seq: ", call.Seq, ", startTime: ", call.StartTime)
	srv.pendingMap.Store(call.Seq, call)
}

// DelPendingCall 删除等待调用
func (srv *App) DelPendingCall(seq uint64) {
	//seelog.Debug("delPendingCall, seq: ", seq)
	srv.pendingMap.Delete(seq)
}

// GetPendingCall 获得等待调用
func (srv *App) GetPendingCall(seq uint64) *idata.PendingCall {
	call, ok := srv.pendingMap.Load(seq)
	if ok {
		return call.(*idata.PendingCall)
	}

	return nil
}

// loopCheckPendingCall 循环检查等待调用
func (srv *App) loopCheckPendingCall(closeSig chan bool) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("loopCheckPendingCall panic:", err, ", Stack: ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-closeSig:
			srv.delTimeoutPendingCall(true)
			return

		case <-ticker.C:
			srv.delTimeoutPendingCall(false)
		}
	}
}

// delTimeoutPendingCall 定时删除等待调用
func (srv *App) delTimeoutPendingCall(force bool) {
	//seelog.Debug("delTimeoutPendingCall")

	srv.pendingMap.Range(
		func(key, value interface{}) bool {
			call := value.(*idata.PendingCall)
			//seelog.Debug("call.StartTime: ", call.StartTime, ", now: ", time.Now().Unix())

			//暂时定为4秒就超时
			if call.StartTime+4 < time.Now().Unix() || force {
				srv.DelPendingCall(key.(uint64))
				retData := &idata.RetData{}
				retData.Err = fmt.Errorf("call timeout")
				call.RetChan <- retData
			}

			return true
		})
}

// GetAppNet 获得app组网
func (srv *App) GetAppNet() iserver.IAppNet {
	return srv.AppNet
}

// startProfServer 开始性能指标检测
func startProfServer() {
	port := viper.GetInt("ServerApp.pprof-port")
	if port == 0 {
		return
	}

	addr := fmt.Sprintf(":%d", port)
	seelog.Infof("Start pprof http://%s/debug/pprof", addr)
	go func() {
		seelog.Info(http.ListenAndServe(addr, nil))
	}()
}
