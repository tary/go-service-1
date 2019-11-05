package servicediscovery

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/spf13/viper"
)

var (
	watchMap        sync.Map
	serviceNewCheck ICheckNewService
)

// ICheckNewService 外部提供用来判断是否是新的服务
type ICheckNewService interface {
	IsNewService(serverID uint64, serverType uint8) bool
}

// IServiceWatcher 关注的服务接口
type IServiceWatcher interface {
	OnNewService(string, uint64)
	GetWatchedServerType() int32
}

// StartDiscovery 开始服务发现
//
// newCheck 提供接口检查一个服务是否是新的服务
// watchers 关注的服务列表
func StartDiscovery(newCheck ICheckNewService, watchers ...IServiceWatcher) error {
	serviceNewCheck = newCheck
	for _, watcher := range watchers {
		watchMap.Store(watcher.GetWatchedServerType(), watcher)
	}

	go loopCheckServer()

	return nil
}

func loopCheckServer() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("loopCheckServer panic:", err, ", Stack: ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-_ctx.Done():
			return

		case <-ticker.C:
			checkNewServer()
		}
	}
}

//checkNewServer 检查是否有新服务器需要连接
//msg 连接成功后发送的消息
//hdl rpc处理函数
func checkNewServer() {
	watchMap.Range(func(k, v interface{}) bool {
		srvType := k.(int32)
		watcher := v.(IServiceWatcher)

		list := getSrvListByType(uint8(srvType))
		for _, v := range list {
			if !serviceNewCheck.IsNewService(v.ServerID, v.Type) {
				watcher.OnNewService(v.SrvOuterAddr, v.ServerID)
			}
		}

		return true
	})
}
