package servicediscovery

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/spf13/viper"
)

//服务加载器
type serverLoader struct {
	serverMapRW sync.RWMutex
	serverMap   map[uint8]ServerList
	interval    time.Duration
	ctx         context.Context
}

// newServerLoader 服务加载器
func newServerLoader(ctx context.Context, t time.Duration) *serverLoader {
	app := &serverLoader{}
	app.interval = t
	app.ctx = ctx
	app.getServerList()

	return app
}

// Start 启动服务加载器
func (app *serverLoader) Start() {
	go app.loop()
}

func (app *serverLoader) loop() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("serverLoader.loop panic:", err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	ticker := time.NewTicker(app.interval)
	defer ticker.Stop()

	for {
		select {
		case <-app.ctx.Done():
			return
		case <-ticker.C:
			app.getServerList()
		}
	}
}

//GetServerByType 获取指定类型的服务器列表
func (app *serverLoader) GetSrvListByType(SrvType uint8) []*ServerInfo {

	app.serverMapRW.RLock()
	defer app.serverMapRW.RUnlock()

	if list, ok := app.serverMap[SrvType]; ok {
		length := list.Len()
		if length <= 0 {
			seelog.Error("Server list is empty, Type %d", SrvType)
			return nil
		}
		return list
	}
	return nil
}

// GetServerList 获取最新的服务器列表
func (app *serverLoader) getServerList() {
	var list []*ServerInfo

	if err := GetServerList(&list); err != nil {
		seelog.Error("redis service is error: ", err)
	}

	app.serverMapRW.Lock()
	defer app.serverMapRW.Unlock()

	app.serverMap = make(map[uint8]ServerList)
	for _, v := range list {
		app.serverMap[v.Type] = append(app.serverMap[v.Type], v)
	}

}
