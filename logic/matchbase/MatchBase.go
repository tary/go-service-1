package matchbase

import (
	"github.com/GA-TECH-SERVER/zeus/logic/matchbase/internal"
	"github.com/GA-TECH-SERVER/zeus/logic/matchbase/matchitf"
)

// MatchBase 匹配基类
type MatchBase struct {
	pm *internal.Manager
}

// Init 初始化
// matchFunction 自定义的匹配接口
// notify 自定义的通知接口
func (mb *MatchBase) Init(matchFunction matchitf.IMatchFunction, notify matchitf.IMatchNotify) {
	mb.pm = internal.NewManager(matchFunction, notify)
}

// TryToMatch 尝试匹配
func (mb *MatchBase) TryToMatch() {
	mb.pm.TryToMatch()
}
