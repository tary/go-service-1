package space

// 实体委托相关代码
// 一个Entity委托给目标Entity后, 该Entity的位置信息可以由目标Entity更新

import (
	"zeus/iserver"
	"zeus/msgdef"
	"zeus/msghandler"
)

type iEntrustCB interface {
	OnEntrust()
}

type iUnEntrustCB interface {
	OnUnEntrust()
}

func (e *Entity) syncEntrustedState(msg *msgdef.SyncUserState) {
	//判断EntityID是否合法, 合法的话转给相应的Entity
	if o, ok := e.entrustedList[msg.EntityID]; ok {
		o.EntrustUpdateSate(msg)
	} else {
		e.Warn("Cant find commission target ID: ", msg.EntityID)
	}
}

// EntrustToTarget 委托给目标Entity
func (e *Entity) EntrustToTarget(target iserver.ISpaceEntity) {
	if target == nil {
		return
	}

	if e.entrustTarget != nil {
		if e.entrustTarget == target {
			return
		}

		e.UnEntrust()
	}

	target.EntrustedBy(e)
	e.entrustTarget = target

	e.setBasePropsDirty(true)

	if iE, ok := e.GetRealPtr().(iEntrustCB); ok {
		iE.OnEntrust()
	}
}

// GetEntrustTarget 获取委托对象
func (e *Entity) GetEntrustTarget() iserver.ISpaceEntity {
	return e.entrustTarget
}

// UnEntrust 取消委托
func (e *Entity) UnEntrust() {
	if e.entrustTarget == nil {
		return
	}

	e.entrustTarget.UnEntrustedBy(e)
	e.entrustTarget = nil

	e.setBasePropsDirty(true)

	if iU, ok := e.GetRealPtr().(iUnEntrustCB); ok {
		iU.OnUnEntrust()
	}
}

// EntrustUpdateSate 处于委托状态时, 接受状态快照
func (e *Entity) EntrustUpdateSate(msg *msgdef.SyncUserState) {
	e.syncClientUserState(msg)
}

// IsEntrusted 是否处于托管状态
func (e *Entity) IsEntrusted() bool {
	return e.entrustTarget != nil
}

// EntrustedBy 被委托回调
func (e *Entity) EntrustedBy(o iserver.ISpaceEntity) {
	if o == nil {
		return
	}

	if e.entrustedList == nil {
		e.entrustedList = make(map[uint64]iserver.ISpaceEntity)
	}
	e.entrustedList[o.GetID()] = o

	e.setBasePropsDirty(true)
}

// UnEntrustedBy 被取消委托时回调
func (e *Entity) UnEntrustedBy(o iserver.ISpaceEntity) {
	if o == nil || e.entrustedList == nil {
		return
	}

	delete(e.entrustedList, o.GetID())

	e.setBasePropsDirty(true)
}

// IsBeEntrusted 是否处于被委托状态
func (e *Entity) IsBeEntrusted() bool {
	return len(e.entrustedList) != 0
}

// EntrustedRPC 调用被委托对象的RPC消息
func (e *Entity) EntrustedRPC(msg *msgdef.RPCMsg) {
	if !e.IsBeEntrusted() {
		e.Warn("EntrustedRPC failed, no entity entrusted to me ", msg)
		return
	}

	target, ok := e.entrustedList[msg.SrcEntityID]
	if !ok {
		e.Warn("EntrustedRPC failed, SrcEntityID error ", msg)
		return
	}

	iM, ok := target.GetRealPtr().(msghandler.IMsgHandlers)
	if !ok {
		e.Warn("EntrustedRPC failed, IMsgHandlers not existed ", msg)
		return
	}

	iM.FireRPC(msg.MethodName, msg.Data)
}
