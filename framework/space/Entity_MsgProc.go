package space

import (
	"zeus/msgdef"
)

// EntityMsgProc Entity消息处理函数
type EntityMsgProc struct {
	e *Entity
}

func (proc *EntityMsgProc) MsgProc_PropsSyncClient(content msgdef.IMsg) {
	msg := content.(*msgdef.PropsSyncClient)
	proc.e.CastMsgToAllClient(msg)
}

func (proc *EntityMsgProc) MsgProc_SyncUserState(content msgdef.IMsg) {
	msg := content.(*msgdef.SyncUserState)
	if msg.EntityID == proc.e.GetID() {
		proc.e.syncClientUserState(msg)
	} else {
		proc.e.syncEntrustedState(msg)
	}
}

func (proc *EntityMsgProc) MsgProc_SessClosed(content interface{}) {
	// log.Info("SessClosed ", proc.e)
	proc.e.SetClient(nil)

	// proc.e.LeaveSpace()
}
