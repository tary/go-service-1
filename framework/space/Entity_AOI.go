package space

import (
	"math"

	"github.com/giant-tech/go-service/base/linmath"
	"github.com/giant-tech/go-service/framework/iserver"
)

const (
	// aoiRange     = 50.0
	aoiTolerance = 1.0
)

type iAOIPacker interface {
	GetID() uint64
	GetType() string
}

// SetWatcher 设置当前entity 为watcher
func (e *Entity) SetWatcher(b bool) {
	e._isWatcher = b
}

// IsWatcher 是否观察者
func (e *Entity) IsWatcher() bool {
	return e._isWatcher
}

func (e *Entity) isBeWatch() bool {
	return e.beWatchedNums > 0
}

func (e *Entity) determineAOIFlag() {
	updataDist := 10.0 //aoiRange * 0.01
	if math.Abs(float64(e.pos.X-e.lastAOIPos.X)) > updataDist ||
		math.Abs(float64(e.pos.Y-e.lastAOIPos.Y)) > updataDist ||
		math.Abs(float64(e.pos.Z-e.lastAOIPos.Z)) > updataDist {
		e.needUpdateAOI = true
	}
}

func (e *Entity) updatePosCoord(pos linmath.Vector3) {

	if e.needUpdateAOI {
		s := e.GetSpace()
		if s != nil {
			s.UpdateCoord(e)
		}

		e.lastAOIPos = pos
		e.needUpdateAOI = false
	}
}

// AddExtWatchEntity 增加额外关注对象
func (e *Entity) AddExtWatchEntity(o iserver.ICoordEntity) {
	if e.extWatchList == nil {
		e.extWatchList = make(map[uint64]*extWatchEntity)
	}

	if _, ok := e.extWatchList[o.GetEntityID()]; ok {
		return
	}

	inMyAOI := false
	if e.GetSpace() != nil {
		e.GetSpace().TravsalAOI(e, func(n iserver.ICoordEntity) {
			// 已经在AOI范围内
			if n.GetEntityID() == o.GetEntityID() {
				inMyAOI = true
			}
		})
	}

	if !inMyAOI {
		e.OnEntityEnterAOI(o)
	}

	e.extWatchList[o.GetEntityID()] = &extWatchEntity{
		entity:  o,
		isInAOI: inMyAOI,
	}
}

// RemoveExtWatchEntity 删除额外关注对象
func (e *Entity) RemoveExtWatchEntity(o iserver.ICoordEntity) {
	if e.extWatchList == nil {
		return
	}

	if _, ok := e.extWatchList[o.GetEntityID()]; !ok {
		return
	}

	inMyAOI := false

	if e.GetSpace() != nil {
		e.GetSpace().TravsalAOI(e, func(n iserver.ICoordEntity) {
			// 已经在AOI范围内
			if n.GetEntityID() == o.GetEntityID() {
				inMyAOI = true
			}
		})
	}

	delete(e.extWatchList, o.GetEntityID())

	if !inMyAOI {
		e.OnEntityLeaveAOI(o)
	}
}

func (e *Entity) clearExtWatchs() {

	for id, we := range e.extWatchList {
		delete(e.extWatchList, id)
		if !we.isInAOI {
			e.OnEntityLeaveAOI(we.entity)
		}
	}
}

// TravsalExtWatchs 遍历额外观察者列表
func (e *Entity) TravsalExtWatchs(f func(*extWatchEntity)) {
	if len(e.extWatchList) == 0 {
		return
	}

	for _, extWatch := range e.extWatchList {
		if !extWatch.isInAOI {
			f(extWatch)
		}
	}
}

//OnEntityEnterAOI 实体进入AOI范围
func (e *Entity) OnEntityEnterAOI(o iserver.ICoordEntity) {
	// 当o在我的额外关注列表中时, 不触发真正的EnterAOI, 只是打个标记
	if extWatch, ok := e.extWatchList[o.GetEntityID()]; ok {
		extWatch.isInAOI = true
		return
	}

	if e._isWatcher {
		e.aoies = append(e.aoies, iserver.AOIInfo{IsEnter: true, Entity: o})
	}

	if o.IsWatcher() {
		e.beWatchedNums++
	}
}

//OnEntityLeaveAOI 实体离开AOI范围
func (e *Entity) OnEntityLeaveAOI(o iserver.ICoordEntity) {
	// 当o在我的额外关注列表中时, 不触发真正的LeaveAOI
	if extWatch, ok := e.extWatchList[o.GetEntityID()]; ok {
		extWatch.isInAOI = false
		return
	}

	if e._isWatcher {
		e.aoies = append(e.aoies, iserver.AOIInfo{IsEnter: false, Entity: o})
	}

	if o.IsWatcher() {
		e.beWatchedNums--
	}
}

func (e *Entity) updateAOI() {
	if len(e.aoies) != 0 && e._isWatcher {
		e.GetRealPtr().(iserver.IAOIUpdate).AOIUpdate(e.aoies)

		e.aoies = e.aoies[0:0]
	}
}

//IsNearAOILayer 是否视野近的层
func (e *Entity) IsNearAOILayer() bool {
	return false
}

//IsAOITrigger 是否要解发AOITrigger事件
func (e *Entity) IsAOITrigger() bool {
	return true //e.IsWatcher()
}
