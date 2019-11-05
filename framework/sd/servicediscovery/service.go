package servicediscovery

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/spf13/viper"
)

var (
	_serverLoader *serverLoader
	_once         sync.Once
)

// GetGroupMgrInst 获取组管理器单例
func startLoader() {
	_once.Do(func() {
		_serverLoader = newServerLoader(_ctx, _getTime)
		_serverLoader.Start()
	})

}

// RegisterService 注册服务
func RegisterService(outerAddr string, sid uint64, stype uint8) error {
	seelog.Debug("RegisterService, outerAddr: ", outerAddr, ", sid: ", sid, ", stype: ", stype)
	svrInfo := &ServerInfo{
		ServerID:     sid,
		Type:         stype,
		SrvOuterAddr: outerAddr,
		Load:         0,
	}

	if outerAddr != "" {
		if err := regState(svrInfo); err != nil {
			return err
		}
	}

	start(svrInfo)

	startLoader()

	return nil
}

//Start 开始协程
func start(svrInfo *ServerInfo) {
	if svrInfo.SrvOuterAddr != "" {
		go refresh(svrInfo)
	}
}

//增加服务器信息
func regState(svrInfo *ServerInfo) error {
	util := newServerUtil(svrInfo.ServerID)
	if util.IsExist() {
		seelog.Error("Server ID is duplicate!!!!! ", svrInfo.ServerID)
		// panic("server is is dupblicate")
	}

	if err := util.Register(svrInfo); err != nil {
		return err
	}
	return nil
}

//Unregister 将服务器信息从redis中删除
func Unregister(svrInfo *ServerInfo) error {
	if err := newServerUtil(svrInfo.ServerID).Delete(); err != nil {
		return err
	}
	return nil
}

// 10秒一次刷新过期时间, 25~35秒一次刷新一次服务器列表
func refresh(svrInfo *ServerInfo) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("refresh panic:", err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ticker := time.NewTicker(_setTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := update(svrInfo); err != nil {
				seelog.Error(err)
			}
		case <-_ctx.Done():
			Unregister(svrInfo)
			return
		}
	}
}

// Update 更新服务器信息
func update(svrInfo *ServerInfo) error {
	return newServerUtil(svrInfo.ServerID).Update(svrInfo)
}

// GetSrvListByType 获取某类型服务器列表
func getSrvListByType(t uint8) []*ServerInfo {
	return _serverLoader.GetSrvListByType(t)
}
