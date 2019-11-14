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
	entity.GroupEntity
	iserver.ICoord
	mapInfo *Map
	mapName string

	tinyEntities map[uint64]ITinyEntity

	startTime time.Time

	isMapLoaded bool

	dirtyEntities []*Entity
}

// OnEntityInit 初始化
func (s *Space) OnEntityInit() {
	s.GroupEntity.OnEntityInit()

	// 暂时先写一个最大的尺寸，后期应该从maploaded结束后再初始化
	s.ICoord = NewTileCoord(9000, 9000)

	s.tinyEntities = make(map[uint64]ITinyEntity)
	s.startTime = time.Now()

	s.isMapLoaded = false
	s.dirtyEntities = make([]*Entity, 0, 100)

	//s.regSpaceSrvID()
}

// OnEntityAfterInit 逻辑层初始化完成之后, 再启动逻辑协程
func (s *Space) OnEntityAfterInit() {
	s.GroupEntity.Entity.OnEntityAfterInit()
	s.loadMap()
}

// OnEntityDestroy 析构函数
func (s *Space) OnEntityDestroy() {
	s.GroupEntity.OnEntityDestroy()
	//s.unRegSpaceSrvID()
}

// func (s *Space) regSpaceSrvID() {

// 	isExistd, err := dbservice.SpaceUtil(s.GetEntityID()).IsExist()

// 	if err != nil {
// 		log.Error("redis error ", err)
// 		return
// 	}

// 	if isExistd {
// 		log.Error("space Id have existed")
// 		return
// 	}

// 	err = dbservice.SpaceUtil(s.GetEntityID()).RegSrvID(iserver.GetSrvInst().GetSrvID())
// 	if err != nil {
// 		log.Error("redis error ", err)
// 	}

// }

// func (s *Space) unRegSpaceSrvID() {
// 	err := dbservice.SpaceUtil(s.GetEntityID()).UnReg()
// 	if err != nil {
// 		log.Error("redis error", err)
// 	}

// }

// OnEntityLoop 完全覆盖Entity的Loop方法
func (s *Space) OnEntityLoop() {
	//s.DoMsg()
	s.GroupEntity.OnEntityLoop()

	s.GroupEntity.Range(func(k, v interface{}) bool {
		if iA, ok := v.(iAOIUpdater); ok {
			iA.updateAOI()
		}

		return true
	})

	// s.refreshEntityState()

	s.GroupEntity.Range(func(k, v interface{}) bool {
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
	//return uint32(time.Now().Sub(s.startTime) / iserver.GetSrvInst().GetFrameDeltaTime())
	return 1
}

// AddEntity 在空间中添加entity
func (s *Space) AddEntity(entityType string, entityID uint64, initParam interface{}, syncInit bool) error {
	e, err := s.CreateEntityWithID(entityType, entityID, s.GetEntityID(), initParam, syncInit, 0)
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
