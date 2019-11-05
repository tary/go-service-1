package matchitf

import (
	"github.com/giant-tech/go-service/logic/matchbase/matchdata"
)

// IMatchFunction 开发者需要实现的的匹配接口
type IMatchFunction interface {
	// MatchForResult 从匹配池中抓取合适的匹配结果
	MatchForResult(pool IMatchPool) []*matchdata.MatchResult

	// MatchForProgress 匹配者是否可以加入等待的房间
	// 匹配者首先找有没有合适的房间，如果有就加入，如果没有就创建一个新房间等待别人加入
	// 如果matcher和room匹配，则返回的第二个参数为true，否则为false
	// 如果此次匹配且room满足开房条件，则返回的第一个参数返回匹配结果，否则返回nil
	MatchForProgress(matcher *matchdata.Matcher, room IMatchPool) (*matchdata.MatchResult, bool)
}

// IMatchNotify 匹配通知
type IMatchNotify interface {
	// MatchFinishNotify 匹配成功通知
	MatchFinishNotify(result *matchdata.MatchResult)
	//玩家进入，用于需要过程的匹配
	MatcherJoinNotify(matcher *matchdata.Matcher, room IMatchPool)
	//玩家退出，用于需要过程的匹配
	MatcherLeaveNotify(matcher *matchdata.Matcher, room IMatchPool)
}

// IMatchPool 匹配池接口
type IMatchPool interface {
	//遍历所有匹配者，f的返回值为false就退出循环
	Range(f func(*matchdata.Matcher) bool)
	// GetCreateTime 获取创建的时间, the number of seconds elapsed since January 1, 1970 UTC.
	GetCreateTime() int64
	// GetMatchMode 获取匹配模式
	GetMatchMode() matchdata.MatchMode
}
