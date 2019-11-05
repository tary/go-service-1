package servermgr

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic" //"github.com/giant-tech/go-service/iserver"

	"github.com/giant-tech/go-service/framework/idata"
	dbservice "github.com/giant-tech/go-service/framework/logicredis"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	serverMapRW sync.RWMutex
	//serverMap   map[uint8]iserver.ServerList
	serverMap  map[uint8]idata.AppList
	succCount  uint32
	totalCount uint32
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer() *LoadBalancer {
	app := &LoadBalancer{}
	if _, err := app.GetServerList(); err != nil {
		log.Error("Init failed", err)
		return nil
	}
	return app
}

// ErrServerBusy 服务器忙
var ErrServerBusy = errors.New("Get server failed, busy")

// GetServerByType 轮询获取服务器
func (app *LoadBalancer) GetServerByType(t uint8) (*idata.AppInfo, error) {
	// 当成功获取服务器的次数超过指定次数时, 重新刷新服务器列表
	if app.succCount >= 50 {
		if _, err := app.GetServerList(); err != nil {
			log.Error("Refresh server list failed", err)
		}
		atomic.StoreUint32(&app.succCount, 0)
	}

	app.serverMapRW.RLock()
	defer app.serverMapRW.RUnlock()

	if list, ok := app.serverMap[t]; ok {
		length := list.Len()
		if length <= 0 {
			return nil, fmt.Errorf("Server list is empty, Type %d", t)
		}

		atomic.AddUint32(&app.totalCount, 1)
		atomic.AddUint32(&app.succCount, 1)
		index := app.totalCount % uint32(length)

		tryTime := 0
		maxLoad := viper.GetInt("Config.MaxLoad")
		if maxLoad == 0 {
			maxLoad = 50
		}
		for {
			if tryTime > length {
				return nil, ErrServerBusy
			}
			tryTime++

			srv := list[index%uint32(length)]
			if srv.Load < maxLoad {
				return srv, nil
			}

			index++
		}
	}
	return nil, fmt.Errorf("Cant get server, Type %d not existed", t)
}

// GetServerList 获取最新的服务器列表
func (app *LoadBalancer) GetServerList() ([]*idata.AppInfo, error) {
	var list []*idata.AppInfo
	if err := dbservice.GetServerList(&list); err != nil {
		return nil, err
	}

	app.serverMapRW.Lock()
	defer app.serverMapRW.Unlock()

	app.serverMap = make(map[uint8]idata.AppList)
	for _, v := range list {
		app.serverMap[v.Type] = append(app.serverMap[v.Type], v)
	}

	atomic.StoreUint32(&app.succCount, 0)
	return list, nil
}

// GetServiceList 获取最新的服务列表
func (app *LoadBalancer) GetServiceList() ([]*idata.ServiceInfo, error) {
	var list []*idata.ServiceInfo
	if err := dbservice.GetServiceList(&list); err != nil {
		return nil, err
	}

	return list, nil
}
