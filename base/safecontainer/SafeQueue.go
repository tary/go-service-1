//使用锁实现的线程安全队列，用来替换safelist
package safecontainer

import "sync"

//SafeQueueNode节点
type SafeQueueNode struct {
	next  *SafeQueueNode
	value interface{}
}
func newNode_M(data interface{}) *SafeQueueNode {
	return &SafeQueueNode{next: nil,value:data}
}

// SafeList 安全链表
type SafeQueue struct {
	head *SafeQueueNode
	tail *SafeQueueNode

	mu sync.Mutex

	C chan bool
}

// NewSafeList 新创建一个列表
func NewSafeList_M() *SafeQueue {
	return &SafeQueue{
		C:make(chan bool, 1),
	}
}

// Put 放入
func (sl *SafeQueue) Put(data interface{}) {
	newNode := newNode_M(data)
	sl.mu.Lock()
	if sl.tail!=nil {
		sl.tail.next = newNode
		sl.tail = newNode
	} else {
		sl.tail = newNode
		sl.head = newNode
	}

	sl.mu.Unlock()

	select {
	case sl.C <- true:
	default:
	}
}

// Pop 拿出
func (sl *SafeQueue) Pop() (interface{}, error) {

	sl.mu.Lock()
	defer sl.mu.Unlock()
	if sl.tail==nil {
		return nil,errNoNode
	}

	if sl.head == sl.tail {
		v := sl.head
		sl.head = nil
		sl.tail = nil
		return v.value,nil
	}

	v := sl.head
	sl.head = sl.head.next
	return v.value,nil
}

// IsEmpty 是否为空
func (sl *SafeQueue) IsEmpty() bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	ret := (sl.tail==nil)
	return ret
}
