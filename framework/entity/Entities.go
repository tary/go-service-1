package entity

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/giant-tech/go-service/base/events"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
	logicredis "github.com/giant-tech/go-service/framework/logicredis"
	"github.com/giant-tech/go-service/framework/net/inet"

	"github.com/giant-tech/go-service/base/timer"
	"github.com/giant-tech/go-service/framework/iserver"

	log "github.com/cihub/seelog"
)

type iEntityCtrl interface {
	//FireMsg(name string, content interface{})
	GetClientSess() inet.ISession
	OnEntityCreated(entityID uint64, entityType string, groupID uint64, protoType interface{}, entities iserver.IEntities, initParam interface{}, syncInit bool, realServerID uint64) error
	OnEntityDestroyed()
	IsDestroyed() bool
	MainLoop()
}

type iEntityRun interface {
	Run()
}

// Entities entity的集合
type Entities struct {
	//msghandler.IRPCHandlers
	ilocal iserver.ILocalService
	*timer.Timer
	*events.GlobalEvents

	isMultilThread bool

	entityCnt int32

	entities       *sync.Map
	entitiesByType *sync.Map
	entitiesByName *sync.Map
}

// NewEntities 创建一个新的Entities
func NewEntities(isMultilThread bool, ilocal iserver.ILocalService) *Entities {
	return &Entities{
		ilocal: ilocal,
		//IRPCHandlers: msghandler.NewRPCHandlers(),
		Timer:        timer.NewTimer(),
		GlobalEvents: events.NewGlobalEventsInst(),

		isMultilThread: isMultilThread,
		entityCnt:      0,

		entities:       &sync.Map{},
		entitiesByType: &sync.Map{},
		entitiesByName: &sync.Map{},
	}
}

// Init 初始化
func (es *Entities) Init() {
	//es.RegMsgProc(es)
}

// Destroy 删除所有的实体
func (es *Entities) Destroy() {

	es.entities.Range(
		func(k, v interface{}) bool {
			if err := es.DestroyEntity(k.(uint64)); err != nil {
				log.Error(err)
			}
			return true
		})

	es.GlobalEvents.Destroy()

	es.entitiesByType = nil
	es.entitiesByName = nil
	es.entities = nil
}

// SyncDestroy 删除所有Entity, 并等待所有Entity删除结束
func (es *Entities) SyncDestroy() {
	es.entities.Range(
		func(k, v interface{}) bool {
			e := es.GetEntity(k.(uint64))
			if err := es.DestroyEntity(k.(uint64)); err != nil {
				log.Error(err)
			}
			for {
				if e.(iEntityState).IsDestroyed() {
					break
				}

				time.Sleep(1 * time.Millisecond)
			}
			return true
		})

	es.GlobalEvents.Destroy()
	es.entitiesByType = nil
	es.entitiesByName = nil
	es.entities = nil
}

// Loop 主动调用loop
func (es *Entities) Loop() {

	//es.DoMsg()

	es.Timer.Loop()
	es.GlobalEvents.HandleEvent()

	if !es.isMultilThread {
		es.entities.Range(func(k, v interface{}) bool {
			v.(iEntityCtrl).MainLoop()
			return true
		})
	}
}

// Range 遍历所有Entity
func (es *Entities) Range(f func(k, v interface{}) bool) {
	es.entities.Range(f)
}

// CreateEntity 创建实体
func (es *Entities) CreateEntity(entityType string, groupID uint64, initParam interface{}, syncInit bool, realServerID uint64) (iserver.IEntity, error) {
	entityID, err := dbservice.GetIDGenerator().GetGlobalID()
	if err != nil {
		return nil, err
	}

	return es.CreateEntityWithID(entityType, entityID, groupID, initParam, syncInit, realServerID)
}

// CreateEntityWithID 创建实体
func (es *Entities) CreateEntityWithID(entityType string, entityID uint64, groupID uint64, initParam interface{}, syncInit bool, realServerID uint64) (iserver.IEntity, error) {
	_, ok := es.entities.Load(entityID)
	if ok {
		return nil, fmt.Errorf("EntityID existed")
	}

	e := es.ilocal.NewEntityByProtoType(entityType)
	if e == nil {
		return nil, fmt.Errorf("Entity type not existed: %s", entityType)
	}

	ie := e.(iEntityCtrl)

	err := ie.OnEntityCreated(entityID, entityType, groupID, ie, es, initParam, syncInit, realServerID)
	if err != nil {
		return nil, err
	}

	if !es.addEntity(entityID, e.(iserver.IEntity)) {
		return nil, fmt.Errorf("addEntity failed, entityType:%s, entityID: %d", entityType, entityID)
	}

	// if es.isMultilThread {
	// 	go func() {
	// 		ticker := time.NewTicker(es.iprotoType.GetTickMS())
	// 		defer ticker.Stop()

	// 		for {
	// 			select {
	// 			case <-ticker.C:
	// 				if ie.IsDestroyed() {
	// 					return
	// 				}
	// 				ie.MainLoop()
	// 			}
	// 		}
	// 	}()
	// }

	if es.isMultilThread {
		irun, ok := e.(iEntityRun)
		if ok {
			go irun.Run()
		}
	}

	atomic.AddInt32(&es.entityCnt, 1)
	return ie.(iserver.IEntity), nil
}

// DestroyEntity 删除Entity
func (es *Entities) DestroyEntity(entityID uint64) error {
	e, ok := es.entities.Load(entityID)
	if !ok {
		return fmt.Errorf("Entity not existed")
	}

	es.delEntity(e.(iserver.IEntity))
	e.(iEntityCtrl).OnEntityDestroyed()

	if !es.isMultilThread {
		e.(iEntityState).OnEntityDestroy()
	}

	atomic.AddInt32(&es.entityCnt, -1)
	return nil
}

// GetEntity 获取Entity
func (es *Entities) GetEntity(entityID uint64) iserver.IEntity {
	if ie, ok := es.entities.Load(entityID); ok {
		return ie.(iserver.IEntity)
	}
	return nil
}

// GetEntityByFunc 获取entity接口
func (es *Entities) GetEntityByFunc(f func(iserver.IEntity) bool) []iserver.IEntity {
	var groupSlice []iserver.IEntity

	es.entities.Range(func(_, v interface{}) bool {
		if f(v.(iserver.IEntity)) {
			groupSlice = append(groupSlice, v.(iserver.IEntity))
		}

		return true
	})

	return groupSlice
}

// TravsalEntity 遍历某一类型的entity
func (es *Entities) TravsalEntity(entityType string, f func(iserver.IEntity)) {
	if it, ok := es.entitiesByType.Load(entityType); ok {
		it.(*sync.Map).Range(func(k, v interface{}) bool {
			//ise := v.(iserver.IEntityStateGetter)
			// if ise.GetEntityState() != iserver.EntityStateLoop {
			// 	return true
			// }

			f(v.(iserver.IEntity))
			return true
		})
	}
}

// addEntity  增加entity
func (es *Entities) addEntity(entityID uint64, e iserver.IEntity) bool {
	_, ok := es.entities.LoadOrStore(entityID, e)
	if ok {
		//已经存在
		log.Error("addEntity, entity exist, entityID: ", entityID)
		e.ClearRPCData()
		return false
	}

	var t *sync.Map
	if it, ok := es.entitiesByType.Load(e.GetType()); ok {
		t = it.(*sync.Map)
	} else {
		t = &sync.Map{}
		es.entitiesByType.Store(e.GetType(), t)
	}

	t.Store(e.GetEntityID(), e)

	return true
}

// delEntity 删除entity
func (es *Entities) delEntity(e iserver.IEntity) {

	es.entities.Delete(e.GetEntityID())
	es.entitiesByName.Delete(e.GetName())

	if it, ok := es.entitiesByType.Load(e.GetType()); ok {
		it.(*sync.Map).Delete(e.GetEntityID())
	}
}

// EntityCount 返回实体数
func (es *Entities) EntityCount() uint32 {
	return uint32(atomic.LoadInt32(&es.entityCnt))
}

// CreateEntityAll 创建实体的所有部分
func (es *Entities) CreateEntityAll(entityType string, initParam interface{}, syncInit bool) (iserver.IEntity, error) {
	e, err := es.CreateEntity(entityType, 0, initParam, syncInit, 0)
	if err != nil {
		return nil, err
	}

	//srvList := global.GetGlobalInst().GetGlobalIntSlice("EntitySrvTypes:" + entityType)
	// srvList, err := dbservice.EntityTypeUtil(entityType).GetSrvType()
	// if err != nil {
	// 	es.DestroyEntity(e.GetID())
	// 	return nil, err
	// }

	//for i := 0; i < len(srvList); i++ {
	// srvType := uint8(srvList[i])
	// if srvType == es.iprotoType.GetServerType() {
	// 	continue
	// }

	/*srvID, err := es.iprotoType.GetSrvIDBySrvType(srvType)
	if err != nil {
		log.Error(err)
		continue
	}
	*/
	//提前注册，这样就可以提前发消息了
	// dbservice.GetEntitySrvUtil(entityID).RegSrvID(srvType, srvID, 0, entityType, dbid)

	// data, err := common.Marshal(initParam)
	// if err != nil {
	// 	log.Error("marshal init error ", err)
	// 	return nil, err
	// }

	/*msg := &msgdef.CreateEntityReq{
		EntityType: entityType,
		EntityID:   entityID,
		CellID:     0,
		InitParam:  serializer.Serialize(initParam),
		DBID:       dbid,
		SrcSrvType: es.iprotoType.GetServerType(),
		SrcSrvID:   es.iprotoType.GetServerID(),
		CallbackID: 0,
	}

	if err := es.iprotoType.PostMsgToCell(srvID, 0, msg); err != nil {
		log.Error(err)
	}*/
	//}

	return e, nil
}

// DestroyEntityAll 销毁实体的所有部分
func (es *Entities) DestroyEntityAll(entityID uint64) error {

	var srvInfos map[uint8]*logicredis.EntitySrvInfo
	var err error

	e := es.GetEntity(entityID)
	if e == nil {
		srvInfos, err = logicredis.GetEntitySrvUtil(entityID).GetSrvIDs()
		if err != nil {
			log.Error("Get entity srv info failed ", err)
			return err
		}
	} else {
		srvInfos = e.GetSrvIDS()
	}

	_ = srvInfos

	// srvInfos, err := dbservice.GetEntitySrvUtil(entityID).GetSrvIDs()
	// if err != nil {
	// 	log.Error("Get entity srv info failed ", err)
	// 	return err
	// }

	return es.DestroyEntity(entityID)
}

//GetEntityByName 根据实体名获取实体
func (es *Entities) GetEntityByName(entityName string) iserver.IEntity {
	if it, ok := es.entitiesByName.Load(entityName); ok {
		return it.(iserver.IEntity)
	}

	return nil
}

//AddEntityByName 根据实体名加入实体
func (es *Entities) AddEntityByName(entityName string, ie iserver.IEntity) error {
	if ie == nil {
		return fmt.Errorf("entity is nil")
	}

	if _, ok := es.entitiesByName.Load(entityName); ok {
		return fmt.Errorf("entity name exist: %s", entityName)
	}

	es.entitiesByName.Store(entityName, ie)

	return nil
}

//DelEntityByName 根据实体名删除实体
func (es *Entities) DelEntityByName(name string) error {
	es.entitiesByName.Delete(name)

	return nil
}

// GetLocalService 获取本地服务接口
func (es *Entities) GetLocalService() iserver.ILocalService {
	return es.ilocal
}

// IsMultiThread 所管理的实体是否为多协程
func (es *Entities) IsMultiThread() bool {
	return es.isMultilThread
}
