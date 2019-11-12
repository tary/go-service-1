package space

import (
	"zeus/iserver"
	"zeus/linmath"
	"zeus/msgdef"
)

type iStateValidate interface {
	StateValidate(oldState, newState IEntityState) int
	StateChange(oldState, newState IEntityState)
}

type iCreateState interface {
	CreateNewEntityState() IEntityState
}

type iStateSender interface {
	SendFullState()
}

// SendFullState 给客户端发送完整状态信息
func (e *Entity) SendFullState() {
	state := e.states.GetLastState().Clone()
	state.SetTimeStamp(e.GetSpace().GetTimeStamp())
	msg := &msgdef.AdjustUserState{
		Data: state.Marshal(),
	}
	e.CastMsgToMe(msg)
}

func (e *Entity) syncClientUserState(msg *msgdef.SyncUserState) {

	stateChecker, ok := e.GetRealPtr().(iStateValidate)

	if !ok {
		panic("no state checker")
	}

	oldState := e.states.GetLastState()

	if oldState.IsModify() {
		return
	}

	var newState IEntityState
	if e.states.cachedState == nil {
		newState = oldState.Clone()
	} else {
		newState = e.states.cachedState
		e.states.cachedState = nil
		oldState.CopyTo(newState)
	}

	newState.Combine(msg.Data)
	var state IEntityState

	// 客户端使用加速挂
	// 暂时只是忽略掉这个状态, 没有额外处理
	if newState.GetTimeStamp() > e.GetSpace().GetTimeStamp() {
		return
	}

	ret := stateChecker.StateValidate(oldState, newState)
	if ret == 2 {
		return
	}
	if ret == 1 {
		state = oldState.Clone()
		state.SetTimeStamp(newState.GetTimeStamp())
		state.SetModify(true)
	} else {
		state = newState
	}

	e.states.addEntityState(state)

	e.SetCoordPos(state.GetPos())

	stateChecker.StateChange(oldState, state)
}

// SetPos 设置位置
func (e *Entity) SetPos(pos linmath.Vector3) {
	if e.IsLinked() {
		return
	}

	e.states.GetLastState().SetPos(pos)
	e.SetCoordPos(pos)
	e.states.GetLastState().SetModify(true)
}

// SetRota 设置旋转
func (e *Entity) SetRota(rota linmath.Vector3) {
	if e.IsLinked() {
		return
	}

	e.states.GetLastState().SetRota(rota)
	e.states.GetLastState().SetModify(true)
}

// GetPos 获取位置
func (e *Entity) GetPos() linmath.Vector3 {
	return e.pos
}

// GetRota 获取旋转
func (e *Entity) GetRota() linmath.Vector3 {
	if e.IsLinked() {
		return e.linkTarget.GetRota()
	}
	return e.states.GetLastState().GetRota()
}

// GetStatePack 获取当前快照的打包
func (e *Entity) GetStatePack() []byte {
	return e.states.GetLastState().Marshal()
}

func (e *Entity) initStates() {
	e.states = NewEntityStates()

	ic, ok := e.GetRealPtr().(iCreateState)
	if !ok {
		panic("the entity must implement createState interface")
	}

	ns := ic.CreateNewEntityState()
	ns.SetTimeStamp(e.GetSpace().GetTimeStamp())
	e.states.addEntityState(ns)
	ns.SetDirty(false)

	ns2 := ic.CreateNewEntityState()
	ns2.SetTimeStamp(e.GetSpace().GetTimeStamp())
	e.states.addEntityState(ns2)
	ns2.SetDirty(false)
}

// GetStates 获取状态快照
func (e *Entity) GetStates() *EntityStates {
	return e.states
}

func (e *Entity) resetState() {
	// if !e.IsMarker() {
	// 	return
	// }

	//如果本身不是观察者，它的状态快照由自己生成
	if !e.IsWatcher() {
		var state IEntityState
		if e.states.cachedState == nil {
			state = e.states.GetLastState().Clone()
		} else {
			state = e.states.cachedState
			e.states.cachedState = nil
			e.states.GetLastState().CopyTo(state)
		}
		state.SetTimeStamp(e.GetSpace().GetTimeStamp())
		e.states.addEntityState(state)
	}
}

func (e *Entity) updateState() {

	// 如果拥有客户端，且状态被服务器修改过，则触发客户端调整消息
	var adjustData []byte
	if e.IsWatcher() {
		ls := e.states.GetLastState()
		if ls.IsModify() {
			msg := &msgdef.AdjustUserState{
				Data: ls.Marshal(),
			}
			e.CastMsgToMe(msg)
			ls.SetModify(false)
			adjustData = msg.Data
		}
	}

	// 如果被观察，把脏数据发给观察者，否则的话，直接修改标记
	//fix me

	if e.isBeWatch() {
		var data []byte
		var isEmpty bool

		// 如果状态被修改过, 则用一个新状态来计算差异, 将差异发送给其他客户端
		if adjustData != nil {
			ic := e.GetRealPtr().(iCreateState)
			ns := ic.CreateNewEntityState()
			ls := e.states.GetLastState()
			ns.SetTimeStamp(ls.GetTimeStamp())
			data, isEmpty = ns.Delta(ls)

			e.states.reflushDirtyState()
		} else {
			data, isEmpty = e.states.reflushDirtyStateAndGetDelta()
		}

		if !isEmpty {
			e.GetSpace().TravsalAOI(e, func(o iserver.ICoordEntity) {
				if e.entrustTarget != nil {
					if e.entrustTarget.GetID() == o.GetID() {
						return
					}
				}

				if o.IsWatcher() && e.GetID() != o.GetID() {
					o.(IWatcher).addStateChangeMsg(e.GetID(), data)
				}
			})

			e.TravsalExtWatchs(func(ext *extWatchEntity) {
				if ext.isInAOI {
					return
				}

				if ext.entity.IsWatcher() {
					ext.entity.(IWatcher).addStateChangeMsg(e.GetID(), data)
				}
			})
		}

	} else {
		e.states.reflushDirtyState()
	}
}

// watcher
func (e *Entity) addStateChangeMsg(id uint64, data []byte) {
	if e.aoiSyncMsg != nil {
		e.aoiSyncMsg.AddData(id, data)
	}
}

func (e *Entity) reflushStateChangeMsg() {
	if e.aoiSyncMsg == nil {
		return
	}

	if e.aoiSyncMsg.Num == 0 {
		return
	}

	if err := e.PostToClient(e.aoiSyncMsg); err != nil {
		e.Error("Send AOISyncMsg failed ", err)
	}

	e.aoiSyncMsg = msgdef.NewAOISyncUserState()
}
