package iserver

import (
	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/idata"
	dbservice "github.com/giant-tech/go-service/framework/logicredis"
	"github.com/giant-tech/go-service/framework/msgdef"
	"github.com/giant-tech/go-service/framework/msghandler"
	"github.com/giant-tech/go-service/framework/net/inet"
)

// IEntity entity接口
type IEntity interface {
	msghandler.IRPCHandlers
	GetEntityID() uint64
	// GetGroupID 获取所在组的ID
	GetGroupID() uint64
	GetType() string
	GetName() string
	SetName(string) error
	GetSrvIDS() map[uint8]*dbservice.EntitySrvInfo
	GetEntitySrvID(srvType uint8) (uint64, uint64, error)
	GetRealPtr() interface{}
	GetIEntities() IEntities

	SetClientSess(inet.ISession) bool
	GetClientSess() inet.ISession
	// Send 发送消息，非协程安全
	Send(msg imsg.IMsg) error

	EnterCell(cellID uint64)
	LeaveCell()

	// PostCallMsg 投递消息给实体
	PostCallMsg(msg *msgdef.CallMsg) error
	// PostCallMsgAndWait 投递消息给实体并等待结果返回
	PostCallMsgAndWait(msg *msgdef.CallMsg) *idata.RetData

	// PostFunction 投递函数给实体，并在实体所在的协程中执行
	PostFunction(f func())
	// PostFunctionAndWait 投递函数给实体协程执行，并等待执行结果
	PostFunctionAndWait(f func() interface{}) interface{}

	SyncCall(sType idata.ServiceType, retPtr interface{}, methodName string, args ...interface{}) error
	AsyncCall(sType idata.ServiceType, methodName string, args ...interface{}) error

	//IsOwnerCellEntity() bool
	//IsCellEntity() bool
}

// IEntityProps Entity属性相关的操作
type IEntityProps interface {
	SetProp(name string, value interface{})
	PropDirty(name string)
	GetProp(name string) interface{}
}

// IEntityPropsSetter Entity属性相关的操作
type IEntityPropsSetter interface {
	SetPropsSetter(IEntityProps)
}

const (
	// EntityStateInit EntityStateInit
	EntityStateInit = 0
	//EntityStateLoop EntityStateLoop
	EntityStateLoop = 1
	//EntityStateDestroy EntityStateDestroy
	EntityStateDestroy = 2
	//EntityStateInValid EntityStateInValid
	EntityStateInValid = 3
)

// IEntityStateGetter 获取Entity状态
type IEntityStateGetter interface {
	GetEntityState() uint8
}
