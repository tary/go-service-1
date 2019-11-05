package servermgr

import (
	"context"
	"runtime/debug"
	"time"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	dbservice "github.com/GA-TECH-SERVER/zeus/framework/logicredis"
)

/*
 负载信息保存在redis中
*/

type iUpdateCtrl interface {
	GetLoad() int
}

// LoadUpdater 负载更新器
type LoadUpdater struct {
	loadCtrl  iUpdateCtrl
	serverID  uint64
	interval  time.Duration
	ctx       context.Context
	ctxCancel context.CancelFunc
}

// NewLoadUpdater 创建负载更新器
// 服务器需要提供GetLoad方法获取负载信息
// 需要指定负载信息更新间隔
func NewLoadUpdater(ctrl iUpdateCtrl, serverID uint64, interval time.Duration) *LoadUpdater {
	app := &LoadUpdater{}
	app.loadCtrl = ctrl
	app.serverID = serverID
	app.interval = interval
	return app
}

// Start 启动负载更新器
func (app *LoadUpdater) Start() {
	app.ctx, app.ctxCancel = context.WithCancel(context.Background())
	go app.loop()
}

// Stop 停止负载更新器
func (app *LoadUpdater) Stop() {
	app.ctxCancel()
}

func (app *LoadUpdater) loop() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("LoadUpdater.loop panic:", err, ", Stack: ", string(debug.Stack()))
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
			app.doUpdateLoad()
		}
	}
}

func (app *LoadUpdater) doUpdateLoad() {
	load := app.loadCtrl.GetLoad()
	if err := dbservice.GetServerUtil(app.serverID).SetLoad(load); err != nil {
		log.Error(err)
	}
}
