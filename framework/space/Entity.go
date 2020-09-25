package space

import (
	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/base/linmath"
	"github.com/giant-tech/go-service/framework/entity"
	"github.com/giant-tech/go-service/framework/iserver"
)

// IEntityCtrl 内部使用接口
type IEntityCtrl interface {
	onEnterSpace()
	onLeaveSpace()
}

// IEnterSpace 上层回调
type IEnterSpace interface {
	OnEnterSpace()
}

// ILeaveSpace 上层回调
type ILeaveSpace interface {
	OnLeaveSpace()
}

// iGameStateSender 发送完整的游戏状态信息
type iGameStateSender interface {
	SendFullGameState()
}

// iLateLooper 内部接口, Entity后处理
type iLateLooper interface {
	onLateLoop()
}

// iAOIUpdater 刷新AOI接口
type iAOIUpdater interface {
	updateAOI()
}

// IWatcher 观察者
type IWatcher interface {
	iserver.IEntity
	iserver.IPos

	PostToClient(imsg.IMsg) error

	/*
		getWatchAOIRange() float32
		markerLeaveAOI(IMarker)
		markerEnterAOI(IMarker)

		addMarker(IMarker)
		removeMarker(IMarker)
		isExistMarker(IMarker) bool
	*/

	addStateChangeMsg(uint64, []byte)
	reflushStateChangeMsg()
}

// // IAOIEnterTrigger 进入aoi区域
// type IAOIEnterTrigger interface {
// 	OnMarkerEnterAOI(m IMarker)
// }

// // IAOILeaveTrigger 离开aoi区域
// type IAOILeaveTrigger interface {
// 	OnMarkerLeaveAOI(m IMarker)
// }

// // IMarker 被观察者
// type IMarker interface {
// 	iserver.IEntity
// 	iserver.IPos
// 	IAOIProp
// 	addWatcher(IWatcher)
// 	removeWatcher(IWatcher)
// 	isExistWatcher(IWatcher) bool
// 	IsMarker() bool
// 	getMarkAOIRange() float32
// }

// Entity 空间中的实体
type Entity struct {
	entity.Entity

	space iserver.ISpace

	rota       linmath.Vector3 // 旋转
	pos        linmath.Vector3 // 当前位置
	lastAOIPos linmath.Vector3 // 上次位置

	needUpdateAOI bool
	_isWatcher    bool

	aoies         []iserver.AOIInfo
	beWatchedNums int

	extWatchList map[uint64]*extWatchEntity
}

// extWatchEntity 额外关注列表
type extWatchEntity struct {
	entity iserver.ICoordEntity

	isInAOI bool // 是否在AOI范围内
}

// OnEntityInit 构造函数
func (e *Entity) OnEntityInit() error {
	e.Entity.OnEntityInit()

	e.pos = linmath.Vector3Invalid()
	e.lastAOIPos = linmath.Vector3Invalid()
	e.needUpdateAOI = false

	e.aoies = make([]iserver.AOIInfo, 0, 5)

	e._isWatcher = false

	return nil
}

// OnEntityAfterInit 后代的初始化
func (e *Entity) OnEntityAfterInit() {
	e.Entity.OnEntityAfterInit()
	e.onEnterSpace()
	e.updatePosCoord(e.pos)
}

// OnEntityDestroy 析构函数
func (e *Entity) OnEntityDestroy() {
	e.onLeaveSpace()

	e.Entity.OnEntityDestroy()
}

// GetSpace 获取所在的空间
func (e *Entity) GetSpace() iserver.ISpace {

	if e.space == nil {
		if e.GetGroupID() == 0 {
			return nil
		}

		e.space = e.GetIEntities().(iserver.ISpace)
	}

	return e.space
}

func (e *Entity) onEnterSpace() {
	ic, ok := e.GetRealPtr().(IEnterSpace)
	if ok {
		ic.OnEnterSpace()
	}

	if e.IsWatcher() {
		// msg := &msgdef.EnterCell{
		// 	CellID: e.GetSpace().GetEntityID(),
		// 	//MapName:  e.GetSpace().GetInitParam().(string),
		// 	EntityID: e.GetEntityID(),
		// 	//Addr:      iserver.GetSrvInst().GetCurSrvInfo().OuterAddress,
		// 	//TimeStamp: e.GetSpace().GetTimeStamp(),
		// }
		// if err := e.Post(iserver.ServerTypeClient, msg); err != nil {
		// 	e.Error("Send EnterSpace failed ", err)
		// }

		e.aoies = append(e.aoies, iserver.AOIInfo{IsEnter: true, Entity: e})
	}
}

func (e *Entity) onLeaveSpace() {
	ic, ok := e.GetRealPtr().(ILeaveSpace)
	if ok {
		ic.OnLeaveSpace()
	}

	if e.GetSpace() != nil {
		e.GetSpace().RemoveFromCoord(e)
	}

	if e.IsWatcher() {
		e.aoies = append(e.aoies, iserver.AOIInfo{IsEnter: false, Entity: e})
		e.clearExtWatchs()
		e.updateAOI()

		//msg := &msgdef.LeaveCell{}
		// if err := e.Post(iserver.ServerTypeClient, msg); err != nil {
		// 	e.Error("Send LeaveSpace failed ", err)
		// }
	}
}

// Entity帧处理顺序
// 处理消息和业务逻辑, 在业务逻辑中会有RPC和CastToAll
// 更新坐标系中的位置
// Space更新所有Entity的AOI
// Space调用所有Entity的LateUpdate, 发送属性消息, 延迟发送消息和缓存的CastToAll消息

// OnLoop 循环调用
func (e *Entity) OnLoop() {

	//e.Entity.DoLooper()
	e.updatePosCoord(e.pos)

	// e.updateAOI()
	// e.resetState()
	// e.Entity.OnLoop()
	// e.updatePosCoord(e.pos)
	// e.updateState()
}

// onLateLoop 后处理
func (e *Entity) onLateLoop() {
	//e.Entity.ReflushDirtyProp()
	//e.Entity.FlushDelayedMsgs()

	// 真正发送所有消息
	//e.FlushDelayedCastMsgs()
}

// IsSpaceEntity 是否是个SpaceEntity
func (e *Entity) IsSpaceEntity() bool {
	return true
}
