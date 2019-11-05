package internal

import (
	"container/list"
	"fmt"
	"time"

	"github.com/GA-TECH-SERVER/zeus/logic/matchbase/matchdata"
	"github.com/GA-TECH-SERVER/zeus/logic/matchbase/matchitf"
)

// NewManager 创建新的Pool Manager
func NewManager(matchFunction matchitf.IMatchFunction, notify matchitf.IMatchNotify) *Manager {
	return &Manager{
		poolMap:       make(map[string]*Pool),
		poolListMap:   make(map[string]*list.List),
		matcherMap:    make(map[string]*MatcherData),
		matchFunction: matchFunction,
		notify:        notify,
	}
}

// MatcherData 匹配者数据
type MatcherData struct {
	pool           *Pool         // 匹配者所在的匹配池
	matcherElement *list.Element // 匹配者在list中的Element
}

// Manager 匹配池管理器
type Manager struct {
	// 匹配函数
	matchFunction matchitf.IMatchFunction
	// 通知接口
	notify matchitf.IMatchNotify
	// [match type, *Pool] 无需知道匹配过程
	poolMap map[string]*Pool
	// [match type, Pool list] 需要追踪匹配过程
	poolListMap map[string]*list.List
	// [matcher key, *MatcherData] 匹配者数据
	matcherMap map[string]*MatcherData
}

// TryToMatch 尝试匹配
func (m *Manager) TryToMatch() {
	for _, value := range m.poolMap {
		results := m.matchFunction.MatchForResult(value)
		for _, result := range results {
			m.matchFinish(value, result)
		}
	}

	var next *list.Element
	for _, value := range m.poolListMap {
		for e := value.Front(); e != nil; e = next {
			next = e.Next()
			result, _ := m.matchFunction.MatchForProgress(nil, e.Value.(matchitf.IMatchPool))
			if result != nil {
				m.matchFinish(e.Value.(*Pool), result)
				//删除临时房间
				value.Remove(e)
			}
		}
	}
}

// InsertMatcher 插入匹配者
func (m *Manager) InsertMatcher(matcher *matchdata.Matcher) error {
	// 如果开始时间为0，则获取当前时间
	if matcher.StartTime == 0 {
		matcher.StartTime = time.Now().Unix()
	}

	err := m.checkMatcher(matcher)
	if err != nil {
		return err
	}

	//如果已经加入匹配池，先从池中删除
	m.RemoveMatcher(matcher.Key, true)

	//加入到匹配池
	if matcher.MatchMode == matchdata.MatchModeResult {
		m.pushMatcherToPool(matcher, m.poolMap[matcher.MatchType])
	} else if matcher.MatchMode == matchdata.MatchModeProgress {
		//匹配者是否有合适的房间加入
		l := m.poolListMap[matcher.MatchType]
		var next *list.Element
		for e := l.Front(); e != nil; e = next {
			next = e.Next()
			result, ok := m.matchFunction.MatchForProgress(matcher, e.Value.(matchitf.IMatchPool))
			if result != nil {
				//成功加入匹配房间
				m.pushMatcherToPool(matcher, e.Value.(*Pool))
				m.matchFinish(e.Value.(*Pool), result)
				//删除临时房间
				l.Remove(e)
				return nil
			}

			if ok {
				//成功加入匹配房间
				m.pushMatcherToPool(matcher, e.Value.(*Pool))
				return nil
			}
		}

		//没有合适的房间，直接创建一个新房间
		newPool := CreateNewPool(matcher.MatchMode)
		m.pushMatcherToPool(matcher, newPool)
		l.PushBack(newPool)
	} else {
		return fmt.Errorf("MatchMode invalid")
	}

	return nil
}

// RemoveMatcher 移除匹配者
func (m *Manager) RemoveMatcher(matcherKey string, notify bool) {
	data, ok := m.matcherMap[matcherKey]
	if !ok {
		//不在匹配池，直接返回
		return
	}

	matcher := data.matcherElement.Value.(*matchdata.Matcher)

	data.pool.Remove(data.matcherElement)

	//是否通知离开
	if notify && data.pool.GetMatchMode() == matchdata.MatchModeProgress {
		m.notify.MatcherLeaveNotify(matcher, data.pool)
	}

	if data.pool.GetMatchMode() == matchdata.MatchModeProgress && data.pool.Len() == 0 {
		//如果是面向过程的匹配，当人数全部退出后需要把pool删除
		poolList := m.poolListMap[matcher.MatchType]
		for e := poolList.Front(); e != nil; e = e.Next() {
			if e.Value.(*Pool) == data.pool {
				poolList.Remove(e)
				break
			}
		}
	}

	delete(m.matcherMap, matcherKey)
}

// checkMatcher 检查匹配者合法性
func (m *Manager) checkMatcher(matcher *matchdata.Matcher) error {
	if matcher.MatchMode == matchdata.MatchModeResult {
		if _, ok := m.poolMap[matcher.MatchType]; !ok {
			m.poolMap[matcher.MatchType] = CreateNewPool(matcher.MatchMode)
		}
	} else if matcher.MatchMode == matchdata.MatchModeProgress {
		if _, ok := m.poolListMap[matcher.MatchType]; !ok {
			m.poolListMap[matcher.MatchType] = list.New()
		}
	}

	return nil
}

// pushMatcherToPool 匹配者加入池中
func (m *Manager) pushMatcherToPool(matcher *matchdata.Matcher, pool *Pool) {
	if pool.GetMatchMode() == matchdata.MatchModeProgress {
		m.notify.MatcherJoinNotify(matcher, pool)
	}

	m.matcherMap[matcher.Key] = &MatcherData{pool: pool, matcherElement: pool.PushBack(matcher)}
}

// matchFinish 匹配结束
func (m *Manager) matchFinish(pool *Pool, result *matchdata.MatchResult) {
	//把匹配完成的匹配者从池子中删除
	for _, matcher := range result.Matchers {
		m.RemoveMatcher(matcher.Key, false)
	}

	m.notify.MatchFinishNotify(result)
}
