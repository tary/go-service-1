package space

import (
	"errors"
	"time"
	"github.com/giant-tech/go-service/framework/entity"
	"github.com/giant-tech/go-service/framework/iserver"
)

// IMapLoader 地图加载成功或失败的接口
type IMapLoader interface {
	OnMapLoadSucceed()
	OnMapLoadFailed()
}

// Space 空间实体
type Space struct {
	entity.Entity
	iserver.ICoord
	*entity.Entities
	mapInfo *Map
	mapName string

	tinyEntities map[uint64]ITinyEntity

	startTime time.Time

	isMapLoaded bool

	dirtyEntities []*Entity
}

// OnInit 初始化
func (s *Space) OnInit() {
	// 暂时先写一个最大的尺寸，后期应该从maploaded结束后再初始化
	s.ICoord = NewTileCoord(9000, 9000)
	s.Entities = entity.NewEntities(false)

	s.tinyEntities = make(map[uint64]ITinyEntity)
	s.startTime = time.Now()

	s.isMapLoaded = false
	s.dirtyEntities = make([]*Entity, 0, 100)

	s.RegMsgProc(&SpaceMsgProc{space: s})
	//s.Entity.IMsgHandlers = s.Entities.IMsgHandlers
	s.Entities.IMsgHandlers = s.Entity.IMsgHandlers

	s.Entity.OnInit()
	//s.regSpaceSrvID()
}

// OnAfterInit 逻辑层初始化完成之后, 再启动逻辑协程
func (s *Space) OnAfterInit() {

	s.Entity.OnAfterInit()
	s.loadMap()
}

// OnDestroy 析构函数
func (s *Space) OnDestroy() {
	s.Entities.Destroy()
	//s.unRegSpaceSrvID()
	s.Entity.OnDestroy()
}

// func (s *Space) regSpaceSrvID() {

// 	isExistd, err := dbservice.SpaceUtil(s.GetID()).IsExist()

// 	if err != nil {
// 		log.Error("redis error ", err)
// 		return
// 	}

// 	if isExistd {
// 		log.Error("space Id have existed")
// 		return
// 	}

// 	err = dbservice.SpaceUtil(s.GetID()).RegSrvID(iserver.GetSrvInst().GetSrvID())
// 	if err != nil {
// 		log.Error("redis error ", err)
// 	}

// }

// func (s *Space) unRegSpaceSrvID() {
// 	err := dbservice.SpaceUtil(s.GetID()).UnReg()
// 	if err != nil {
// 		log.Error("redis error", err)
// 	}

// }

// OnLoop 完全覆盖Entity的Loop方法
func (s *Space) OnLoop() {
	//s.DoMsg()
	s.Entity.OnLoop()
	s.Entities.MainLoop()

	s.Entities.Range(func(k, v interface{}) bool {
		if iA, ok := v.(iAOIUpdater); ok {
			iA.updateAOI()
		}

		return true
	})

	// s.refreshEntityState()

	s.Entities.Range(func(k, v interface{}) bool {
		if iW, ok := v.(IWatcher); ok {
			if iW.GetType() == "Player" {
				iW.reflushStateChangeMsg()
			}
		}

		if iL, ok := v.(iLateLooper); ok {
			iL.onLateLoop()
		}

		return true
	})
}

// GetTimeStamp 获取当前的时间戳
func (s *Space) GetTimeStamp() uint32 {
	return uint32(time.Now().Sub(s.startTime) / iserver.GetSrvInst().GetFrameDeltaTime())
}

// AddEntity 在空间中添加entity
func (s *Space) AddEntity(entityType string, entityID uint64, dbid uint64, initParam interface{}, syncInit bool, isGhost bool) error {
	e, err := s.CreateEntity(entityType, entityID, dbid, s.GetID(), initParam, syncInit, isGhost, 0)
	if err != nil {
		return err
	}

	_, ok := e.(iserver.ICellEntity)
	if !ok {
		s.DestroyEntity(e.GetEntityID())
		return errors.New("the entity which add to space must be ICellEntity ")
	}

	return nil
}

// RemoveEntity 在空间中删除entity
func (s *Space) RemoveEntity(entityID uint64) error {

	e := s.GetEntity(entityID)
	if e == nil {
		return errors.New("no entity exist")
	}

	return s.DestroyEntity(entityID)
}

/*
	Space 本身是一个实体，但是消息相关的几个函数都由Entities来执行
	主要原因是为了方便的使用几个Entities的默认处理函数
	如Entity的消息传送等等
*/

// IsSpace 是否是空间
func (s *Space) IsSpace() bool {
	return true
}

// MainLoop 主循环
func (s *Space) MainLoop() {
	s.Entity.MainLoop()
}
