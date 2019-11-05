package matchbase

import (
	"github.com/cihub/seelog"
	"github.com/GA-TECH-SERVER/zeus/logic/matchbase/matchdata"
)

// RPCMatchReq 请求匹配，返回错误内容，如果为空串说明没有错误
func (mb *MatchBase) RPCMatchReq(matcher *matchdata.Matcher) string {
	seelog.Debug("RPCMatchReq, key: ", matcher.Key, ", mode: ", matcher.MatchMode)

	err := mb.pm.InsertMatcher(matcher)
	if err != nil {
		seelog.Error("RPCMatchReq, matchKey: ", matcher.Key, ", err: ", err)
		return err.Error()
	}

	return ""
}

// RPCCancleMatchReq 取消匹配
func (mb *MatchBase) RPCCancleMatchReq(matcherKey string) {
	seelog.Debug("RPCCancleMatchReq, key: ", matcherKey)

	mb.pm.RemoveMatcher(matcherKey, true)
}
