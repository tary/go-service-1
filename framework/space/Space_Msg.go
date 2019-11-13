package space

import (
	"github.com/giant-tech/go-service/framework/iserver"
)

type iSpaceCtrl interface {
	updateCoord(iserver.ICellEntity)
}

func (s *Space) refreshEntityState() {

	s.TravsalEntity("Player", func(entity iserver.IEntity) {
		iw, ok := entity.(IWatcher)
		if ok {
			iw.reflushStateChangeMsg()
		}
	})
}

// SetEnable 临时处理
func (s *Space) SetEnable(enable bool) {
	s.Entity.SetEnable(enable)
}

// FireMsg 触发消息
func (s *Space) FireMsg(name string, content interface{}) {
	s.Entity.FireMsg(name, content)
}

// FireRPC 触发RPC消息
func (s *Space) FireRPC(methodName string, data []byte) {
	s.Entity.FireRPC(methodName, data)
}

// RegMsgProc 注册消息处理函数
func (s *Space) RegMsgProc(proc interface{}) {
	s.Entity.RegMsgProc(proc)
}

// DoNormalMsg 立刻处理消息
func (s *Space) DoNormalMsg(name string, data interface{}) error {
	return s.Entity.DoNormalMsg(name, data)
}

// DoRPCMsg 立刻处理RPC消息
func (s *Space) DoRPCMsg(name string, data []byte) error {
	return s.Entity.DoRPCMsg(name, data)
}
