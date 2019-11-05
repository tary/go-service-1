package iserver

//"github.com/giant-tech/go-service/msghandler_entity"

// IEntities 用于实体的管理
type IEntities interface {
	//msghandler.IRPCHandlers
	//timer.ITimer
	//events.IEvents

	//销毁所有实体
	Destroy()
	//每帧调用
	Loop()

	CreateEntityAll(entityType string, initParam interface{}, syncInit bool) (IEntity, error)
	DestroyEntityAll(entityID uint64) error

	CreateEntity(entityType string, groupID uint64, initParam interface{}, syncInit bool, realServerID uint64) (IEntity, error)
	CreateEntityWithID(entityType string, entityID uint64, groupID uint64, initParam interface{}, syncInit bool, realServerID uint64) (IEntity, error)

	DestroyEntity(entityID uint64) error

	GetEntity(entityID uint64) IEntity
	GetEntityByFunc(f func(IEntity) bool) []IEntity

	//GetEntityByName(entityName string) IEntity
	AddEntityByName(entityName string, ie IEntity) error
	//仅仅从名字管理器中移除
	DelEntityByName(entityName string) error

	TravsalEntity(entityType string, f func(IEntity))

	EntityCount() uint32
	GetLocalService() ILocalService

	// IsMultiThread 所管理的实体是否为多协程
	IsMultiThread() bool
}
