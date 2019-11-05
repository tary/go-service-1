package servicediscovery

import (
	"context"
	"time"
)

var (
	_setTime   = 10 * time.Second
	_getTime   = 5 * time.Second
	_ctx       context.Context
	_ctxCancel context.CancelFunc
)

func init() {
	_ctx, _ctxCancel = context.WithCancel(context.Background())
	_ = _ctxCancel
}

// SetInterval 设置服务发现内部的时间间隔
// setTime: 向redis注册自己的时间间隔
// getTime: 从redis读取所有服务的时间间隔
func SetInterval(setTime, getTime time.Duration) {
	_setTime, _getTime = setTime, getTime
}

// Stop 停止服务发现模块
func Stop() {
	_ctxCancel()
}
