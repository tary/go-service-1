package dailytimer

import (
	"container/list"
	"reflect"
	"time"

	log "github.com/cihub/seelog"
)

//type Job func(TaskData)

//TaskData 回调函数的参数类型
type TaskData map[interface{}]interface{}

//TimeWheel 时间轮
type TimeWheel struct {
	location          *time.Location //时区：默认是Asia/Shanghai
	interval          time.Duration  //时间轮转动一格的时间，单位可选：ns,us,ms,s,m,h，默认1s
	ticker            *time.Ticker
	slots             []*list.List          //时间轮槽
	timer             map[interface{}]int32 //key 定时器任务唯一标识，value 定时器指向的槽
	currentPos        int32                 //当前指针指向哪一个槽
	slotNum           int32                 //槽数量
	addTaskChannel    chan Task             //新增任务channel
	removeTaskChannel chan interface{}      //删除任务channel
	stopChannel       chan bool             //停止定时器channel
}

//Task 任务
type Task struct {
	//startTime string          //定时任务开启时间,默认格式："2018-12-27 14:48:00"
	flag      bool            //true为重复执行任务，false为单次执行任务
	delay     time.Duration   //每隔多久执行一次，为0时表示只执行一次
	circle    int32           //时间轮需要转动几圈
	key       interface{}     //定时器唯一标识，用于删除定时器
	data      []reflect.Value //回调函数参数
	timerFunc reflect.Value   //回调函数
}

//New 新建一个时间轮
func New(d string, slotNum int32) *TimeWheel {
	interval, err := time.ParseDuration(d)
	if err != nil {
		log.Error(err)
	}
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[interface{}]int32),
		slotNum:           slotNum,
		currentPos:        0,
		addTaskChannel:    make(chan Task, 100),
		removeTaskChannel: make(chan interface{}, 100),
		stopChannel:       make(chan bool),
	}
	tw.location, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Error(err)
	}
	tw.initSlots()
	return tw
}

//初始化槽，每个槽指向一个双向链表
func (tw *TimeWheel) initSlots() {
	for i := int32(0); i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

//Start 启动时间轮
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

//Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

//AddTimerTask 添加定时器任务,startTime为任务开始时间，delay为任务执行周期，为0s时只执行一次， key作为定时器任务唯一标识
func (tw *TimeWheel) AddTimerTask(startTime string, delay time.Duration, key interface{}, flag bool, jobFunc interface{}, jobArgs ...interface{}) {
	TaskValue := reflect.ValueOf(jobFunc)
	if TaskValue.Kind() != reflect.Func {
		log.Error("only function can be schedule.")
	}
	if len(jobArgs) != TaskValue.Type().NumIn() {
		log.Error("The number of args valid.")
	}
	in := make([]reflect.Value, len(jobArgs))
	for i, arg := range jobArgs {
		in[i] = reflect.ValueOf(arg)
	}
	task := Task{
		flag:      flag,
		data:      in,
		key:       key,
		delay:     delay,
		timerFunc: TaskValue,
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, tw.location)
	if err != nil {
		log.Debug(err)
	}
	now := time.Now()
	t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, tw.location)
	var delayTime time.Duration
	if now.After(t) {
		delayTime = 0
	} else {
		delayTime = t.Sub(now)
	}
	go func(d time.Duration) {
		if delay <= 0 || key == nil {
			return
		}
		time.Sleep(d)
		if d > 0 {
			go task.timerFunc.Call(task.data)
		}
		tw.addTaskChannel <- task
	}(delayTime)

}

//RemoveTimer 删除定时器 key为添加定时器时传递定时器的唯一标识
func (tw *TimeWheel) RemoveTimer(key interface{}) {
	if key == nil {
		return
	}
	tw.removeTaskChannel <- key
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

// scanAndRunTask 扫描链表中过期的定时器，并执行回调函数
func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}
		go task.timerFunc.Call(task.data)
		if task.flag == true {
			tw.updateTask(task)
		}
		next := e.Next()
		l.Remove(e)
		delete(tw.timer, task.key)
		e = next
	}
}

// addTask 新增任务到链表中
func (tw *TimeWheel) addTask(task *Task) {
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	tw.slots[pos].PushBack(task)
	tw.timer[task.key] = pos
}

//getPositionAndCircle 获取定时器在槽中的位置，时间轮需要转动的圈数
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int32, circle int32) {
	delaySeconds := int32(d.Seconds())
	intervalSeconds := int32(tw.interval.Seconds())
	circle = int32(delaySeconds / intervalSeconds / tw.slotNum)
	pos = int32(tw.currentPos-1+delaySeconds/intervalSeconds) % tw.slotNum
	return pos, circle
}

func (tw *TimeWheel) updateTask(task *Task) {
	tw.addTaskChannel <- *task
}

//从链表中删除任务
func (tw *TimeWheel) removeTask(key interface{}) {
	//获取定时器所在的槽
	position, ok := tw.timer[key]
	if !ok {
		return
	}
	//获取槽指向的链表
	l := tw.slots[position]
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.key == key {
			delete(tw.timer, task.key)
			l.Remove(e)
		}
		e = e.Next()
	}
}
