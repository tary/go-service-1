package internal

import (
	"container/list"
	"time"

	"github.com/giant-tech/go-service/logic/matchbase/matchdata"
)

// CreateNewPool 创建新Pool
func CreateNewPool(mode matchdata.MatchMode) *Pool {
	return &Pool{List: list.New(), matchMode: mode, createTime: time.Now().Unix()}
}

// Pool 匹配池
type Pool struct {
	*list.List
	matchMode  matchdata.MatchMode // 匹配模式
	createTime int64               // 创建的时间
}

// Range 遍历
func (p *Pool) Range(f func(*matchdata.Matcher) bool) {
	for e := p.Front(); e != nil; e = e.Next() {
		if !f(e.Value.(*matchdata.Matcher)) {
			break
		}
	}
}

// GetCreateTime 获取创建的时间, the number of seconds elapsed since January 1, 1970 UTC.
func (p *Pool) GetCreateTime() int64 {
	return p.createTime
}

// GetMatchMode 获取匹配模式
func (p *Pool) GetMatchMode() matchdata.MatchMode {
	return p.matchMode
}
