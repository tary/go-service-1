package dailytimer

import (
	"sync"
)

//TimerMgr 定时器管理器变量
var TimerMgr *TimerManager

//TimerManager 定时器管理器
type TimerManager struct {
	TimeMap *sync.Map
}

//InitTimerManager 初始化定时器管理器
func InitTimerManager() {
	TimerMgr = &TimerManager{}
	TimerMgr.TimeMap = &sync.Map{} //key: string,value TimerWheel
}
