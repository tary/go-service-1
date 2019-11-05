package entity

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/base/serializer"
	"github.com/GA-TECH-SERVER/zeus/framework/idata"
	"github.com/GA-TECH-SERVER/zeus/framework/iserver"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"
	"github.com/GA-TECH-SERVER/zeus/framework/msghandler"

	dbservice "github.com/GA-TECH-SERVER/zeus/framework/logicredis"

	"github.com/GA-TECH-SERVER/zeus/framework/servicedef"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// GameDBName game db name
//var GameDBName = "game"

// iEntityState Entity的状态相关函数
type iEntityState interface {
	OnEntityInit()
	OnEntityAfterInit() error
	OnEntityLoop()
	OnEntityDestroy()
	IsDestroyed() bool
}

// iEntityState 留给后代作一些回调方法
type iEntityInit interface {
	OnInit(interface{}) error
}

type iEntityLoop interface {
	OnLoop()
}

type iEntityDestroy interface {
	OnDestroy()
}

// Entity 代表服务器端一个可通讯对象
type Entity struct {
	msghandler.RPCHandlers
	entityType string
	entityID   uint64
	entityName string
	groupID    uint64 //实体所在group的ID,如果groupID不为0，意味着此实体属于某一组织，为队伍成员，房间成员等
	//groupType  uint32 //实体所在group的type
	CliSess inet.ISession

	//
	realPtr     interface{}
	entitiesPtr iserver.IEntities
	state       byte
	ieState     iEntityState
	initParam   interface{}

	// props begin
	props          map[string]*PropInfo
	def            *Def
	dirtyPropList  []*PropInfo
	dirtySaveProps map[string]*PropInfo
	dirtySyncProps map[uint32][]*PropInfo
	// props end

	srvIDS    map[uint8]*dbservice.EntitySrvInfo
	srvIDSMux *sync.RWMutex

	// DataC 是向本实体发送消息的Channel
	DataC chan *idata.CallData
}

// DBData entity初始db数据
/*type DBData struct {
	ID   bson.ObjectId `bson:"_id"`
	DBID uint64        `bson:"dbid"`
}

// NewDBData 创建带有默认值的数据
func NewDBData() *DBData {
	return &DBData{}
}

// InsertData 往数据库插入entity初始数据
func InsertData(p *DBData, entitytype string) {

	if p == nil {
		return
	}

	p.ID = bson.NewObjectId()
	mongodbservice.MongoDBInsert(GameDBName, entitytype, p)
}
*/

// SetClientSess 设置客户端session
func (e *Entity) SetClientSess(sess inet.ISession) bool {
	// 两个都不为nil说明有问题
	if e.CliSess != nil && sess != nil {
		log.Error("cliSess is not nil")
		return false
	}

	e.CliSess = sess

	return true
}

// GetInitParam 获取初始化参数
func (e *Entity) GetInitParam() interface{} {
	return e.initParam
}

// GetClientSess 获取客户端session
func (e *Entity) GetClientSess() inet.ISession {
	if e.CliSess != nil {
		return e.CliSess
	}

	return nil
}

// GetName 获取实体Name
func (e *Entity) GetName() string {
	return e.entityName
}

// GetEntityID 获取实体ID
func (e *Entity) GetEntityID() uint64 {
	return e.entityID
}

// SetEntityID 设置实体ID， 篮球临时接口
func (e *Entity) SetEntityID(eid uint64) {
	e.entityID = eid
}

// GetGroupID 获取实体groupID
func (e *Entity) GetGroupID() uint64 {
	return e.groupID
}

// GetGroupType 获取实体groupType
// func (e *Entity) GetGroupType() uint32 {
// 	return e.groupType
// }

// GetType 获取实体类型
func (e *Entity) GetType() string {
	return e.entityType
}

// GetRealPtr 获取真实的后代对象的指针
func (e *Entity) GetRealPtr() interface{} {
	return e.realPtr
}

// GetSrvIDS 获取玩家的分布式实体所在的服务器列表
func (e *Entity) GetSrvIDS() map[uint8]*dbservice.EntitySrvInfo {
	return e.srvIDS
}

// GetEntitySrvID 获取实体的服务ID
func (e *Entity) GetEntitySrvID(srvType uint8) (uint64, uint64, error) {

	e.srvIDSMux.RLock()
	srvID, ok := e.srvIDS[srvType]
	e.srvIDSMux.RUnlock()

	if ok {
		return srvID.SrvID, srvID.GroupID, nil
	}

	// 第一次尝试不成功, 则先刷新一次信息
	e.RefreshSrvIDS()

	e.srvIDSMux.RLock()
	srvID, ok = e.srvIDS[srvType]
	e.srvIDSMux.RUnlock()

	if !ok {
		return 0, 0, fmt.Errorf("Entity srvType [%d] not existed", srvType)
	}

	return srvID.SrvID, srvID.GroupID, nil
}

// SetName 设置实体Name并更新name管理器
func (e *Entity) SetName(name string) error {
	if name == e.entityName {
		return nil
	}

	if len(name) > 0 {
		ie := e.entitiesPtr.AddEntityByName(name, e.realPtr.(iserver.IEntity))
		if ie != nil {
			//已经存在
			log.Error("entity name exist: ", name)
			return fmt.Errorf("entity name exist: %s", name)
		}
	}

	if len(e.entityName) > 0 {
		e.entitiesPtr.DelEntityByName(e.entityName)
	}

	e.entityName = name

	return nil
}

// OnEntityCreated 初始化
func (e *Entity) OnEntityCreated(entityID uint64, entityType string, groupID uint64, protoType interface{}, entities iserver.IEntities, initParam interface{}, syncInit bool, realServerID uint64) error {

	//e.IRPCHandlers = msghandler.NewRPCHandlers()
	//只有多线程的时候才需要，如果是单协程就直接交给
	if entities.IsMultiThread() {
		e.DataC = make(chan *idata.CallData, 10240)
	}

	e.entityType = entityType
	e.entityID = entityID

	e.groupID = groupID
	e.initParam = initParam

	//e.realServerID = realServerID

	e.realPtr = protoType
	e.entitiesPtr = entities

	e.state = iserver.EntityStateInit
	e.ieState = protoType.(iEntityState)

	e.srvIDS = make(map[uint8]*dbservice.EntitySrvInfo)
	e.srvIDSMux = &sync.RWMutex{}

	e.props = make(map[string]*PropInfo)
	e.dirtyPropList = make([]*PropInfo, 0, 1)
	e.dirtySaveProps = make(map[string]*PropInfo)
	e.dirtySyncProps = make(map[uint32][]*PropInfo)

	ps, ok := e.realPtr.(iserver.IEntityPropsSetter)
	if ok {
		ps.SetPropsSetter(e.realPtr.(iserver.IEntityProps))
	}

	var err error

	if syncInit {
		ies := e.ieState
		ies.OnEntityInit()
		err = ies.OnEntityAfterInit()
	}

	return err
}

// OnEntityRegSrvID 注册serverID
func (e *Entity) OnEntityRegSrvID() {
	e.RegSrvID()
}

// OnEntityDestroyed 当Entity销毁时调用
func (e *Entity) OnEntityDestroyed() {
	if e.state == iserver.EntityStateLoop || e.state == iserver.EntityStateInit {
		e.state = iserver.EntityStateDestroy
		e.SavePropsToDB()
		e.UnregSrvID()

		ii, ok := e.GetRealPtr().(iEntityDestroy)
		if ok {
			ii.OnDestroy()
		}

	}
}

// MainLoop 主循环
func (e *Entity) MainLoop() {

	defer func() {
		if err := recover(); err != nil {
			log.Error(err, e)
			if viper.GetString("Config.Recover") == "0" {
				panic(fmt.Sprintln(err, e))
			}
		}
	}()

	ies := e.ieState

	switch e.state {
	case iserver.EntityStateInit:
		{
			ies.OnEntityInit()
			ies.OnEntityAfterInit()
		}
	case iserver.EntityStateLoop:
		{
			ies.OnEntityLoop()
		}
	case iserver.EntityStateDestroy:
		{
			ies.OnEntityDestroy()
			e.state = iserver.EntityStateInValid
		}
	default:
		{
			// do nothing
		}
	}
}

// OnEntityInit entity init
func (e *Entity) OnEntityInit() {
	e.state = iserver.EntityStateLoop

	// 初始创建entity数据
	/*type DBMap bson.M
	var tempDBMap DBMap

	if len(tempDBMap) == 0 {
		dbdata := NewDBData()
		InsertData(dbdata, e.GetName())
	}
	*/
	e.InitProp(GetDefs().GetDef(e.entityType), true)

	e.RegSrvID()

	e.RegRPCMsg(e.realPtr)
}

// OnEntityAfterInit 实体创建之后的初始化
func (e *Entity) OnEntityAfterInit() error {
	var err error

	ii, ok := e.GetRealPtr().(iEntityInit)
	if ok {
		err = ii.OnInit(e.GetInitParam())
		//校验实体rpc方法
		rpchandlers, err2 := e.GetRPCHandlers()
		if err2 == nil {
			def := servicedef.GetServiceDefs().GetDef("Entity")
			if def != nil {
				for methodname := range def.Methods {
					_, ok := rpchandlers.Load(methodname)
					if !ok {
						log.Error("name = Entity", " method: ", methodname, " not implement")
					}
				}
			}
		}
	} else {
		log.Error("the entity ", e.GetType(), " no init method")
		err = fmt.Errorf("entity didnot have init method")
	}

	return err
}

// OnEntityLoop entity循环
func (e *Entity) OnEntityLoop() {
	ii, ok := e.GetRealPtr().(iEntityLoop)
	if ok {
		ii.OnLoop()
	}
}

// OnEntityDestroy entity销毁
func (e *Entity) OnEntityDestroy() {

}

// IsDestroyed IsDestroyed是否销毁
func (e *Entity) IsDestroyed() bool {
	return e.state == iserver.EntityStateInValid
}

// PostCallMsg 把消息投递给实体，立即返回
func (e *Entity) PostCallMsg(msg *msgdef.CallMsg) error {
	if !e.GetIEntities().IsMultiThread() {
		panic("Not MultiThread ")
	}

	data := &idata.CallData{}
	data.Msg = msg

	e.DataC <- data

	return nil
}

// PostCallMsgAndWait 把消息投递给实体并且等待返回结果
func (e *Entity) PostCallMsgAndWait(msg *msgdef.CallMsg) *idata.RetData {
	if !e.GetIEntities().IsMultiThread() {
		panic("Not MultiThread ")
	}

	data := &idata.CallData{}
	data.Msg = msg

	// 结果从ChanRet返回
	data.ChanRet = make(chan *idata.RetData, 1)
	e.DataC <- data

	// 等待直到返回结果
	return <-data.ChanRet

}

// PostFunction 向user协程投递函数
func (e *Entity) PostFunction(f func()) {
	if !e.GetIEntities().IsMultiThread() {
		panic("Not MultiThread ")
	}

	e.DataC <- &idata.CallData{Func: f}
}

// PostFunctionAndWait calls a function f and returns the result.
// f runs in the LobbyUser's goroutine.
func (e *Entity) PostFunctionAndWait(f func() interface{}) interface{} {
	if !e.GetIEntities().IsMultiThread() {
		panic("Not MultiThread ")
	}

	// 结果从ch返回
	ch := make(chan interface{}, 1)
	e.DataC <- &idata.CallData{Func: func() { ch <- f() }}

	// 等待直到返回结果
	return <-ch
}

// HandleCallMsg 处理调用消息
func (e *Entity) HandleCallMsg() {
	if !e.GetIEntities().IsMultiThread() {
		panic("Not MultiThread ")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Error("Entity.HandleCallMsg panic:", err, ", ", string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	maxCount := 1000
	for i := 0; i < maxCount; i++ {
		select {
		case data := <-e.DataC:
			e.ProcessCall(data)

		default:
			return
		}
	}
}

//ProcessCall 处理单个调用
func (e *Entity) ProcessCall(data *idata.CallData) {
	//如果Func不为nil则直接调用
	if data.Func != nil {
		data.Func()
		return
	}

	e.DoRPCMsg(data.Msg.MethodName, data.Msg.Params, data.ChanRet)
}

// GetIEntities 获取实体管理器
func (e *Entity) GetIEntities() iserver.IEntities {
	return e.entitiesPtr
}

// IsGhost 是否为ghost
func (e *Entity) IsGhost() bool {
	return false
}

// SyncCall 同步调用，等待返回
func (e *Entity) SyncCall(sType idata.ServiceType, retPtr interface{}, methodName string, args ...interface{}) error {
	//目前不允许调用自己
	if e.GetIEntities().GetLocalService().GetSType() == sType {
		return fmt.Errorf("Can't call self")
	}

	sid, gID, err := e.GetEntitySrvID(uint8(sType))
	if err != nil {
		return err
	}

	s := iserver.GetServiceProxyMgr().GetServiceByID(sid)
	if !s.IsValid() {
		return fmt.Errorf("service not exist: %d", sid)
	}

	msg := &msgdef.CallMsg{}

	msg.FromSID = e.GetIEntities().GetLocalService().GetSID()
	msg.SType = uint8(sType)
	msg.SID = sid
	msg.GroupID = gID
	msg.MethodName = methodName
	msg.IsSync = true
	msg.EntityID = e.entityID
	msg.Params = serializer.SerializeNew(args...)

	//目前不支持对客户端的同步调用
	if sType == idata.ServiceClient {
		log.Errorf("SyncCall failed: can't syncCall client method")
		return fmt.Errorf("SyncCall failed: can't syncCall client method")
	}

	// //发送给客户端
	// if e.GetIEntities().GetLocalService().GetSType() == idata.ServiceGateway && sType == idata.ServiceClient {
	// 	cli := e.GetClientSess()
	// 	if cli != nil {
	// 		msg.Seq = iserver.GetApp().GetSeq()
	// 		cli.Send(msg)

	// 		//加入到pending中
	// 		call := &idata.PendingCall{}
	// 		call.RetChan = make(chan *idata.RetData, 1)
	// 		call.Seq = msg.Seq
	// 		call.MethodName = methodName
	// 		call.Reply = retPtr
	// 		call.ToServiceID = sid

	// 		iserver.GetApp().AddPendingCall(call)

	// 		retData := <-call.RetChan
	// 		if retData.Err != nil {
	// 			return retData.Err
	// 		}

	//		if retPtr != nil {
	//		serializer.UnSerializeNew(retPtr, retData.Ret)
	//		}

	// 		return nil
	// 	} else {
	// 		log.Errorf("SyncCall failed: ClientSess is null, sType: %d", sType)
	// 		return fmt.Errorf("SyncCall failed: ClientSess is null, sType: %d", sType)
	// 	}
	// }

	if s.IsLocal() {
		//直接发送
		is := iserver.GetLocalServiceMgr().GetLocalService(sid)
		if is == nil {
			log.Error("SyncCall failed, service not found: ", sid)
			return fmt.Errorf("SyncCall failed, service not found: %d", sid)
		}

		retData := is.PostCallMsgAndWait(msg)
		if retData.Err != nil {
			return retData.Err
		}

		if retPtr != nil {
			if err := serializer.UnSerializeNew(retPtr, retData.Ret); err != nil {
				return err
			}
		}
	} else {
		msg.Seq = iserver.GetApp().GetSeq()
		s.SendMsg(msg)

		//加入到pending中
		call := &idata.PendingCall{}
		call.RetChan = make(chan *idata.RetData, 1)
		call.Seq = msg.Seq
		call.MethodName = methodName
		call.Reply = retPtr
		call.ToServiceID = sid

		iserver.GetApp().AddPendingCall(call)

		retData := <-call.RetChan
		if retData.Err != nil {
			return retData.Err
		}

		if retPtr != nil {
			if err := serializer.UnSerializeNew(retPtr, retData.Ret); err != nil {
				return err
			}
		}
	}

	return nil
}

// AsyncCall 异步调用，立即返回
func (e *Entity) AsyncCall(sType idata.ServiceType, methodName string, args ...interface{}) error {
	//目前不允许调用自己
	if e.GetIEntities().GetLocalService().GetSType() == sType {
		log.Error("AsyncCall Can't call self, methodName: ", methodName)
		return fmt.Errorf("Can't call self")
	}

	sid, gID, err := e.GetEntitySrvID(uint8(sType))
	if err != nil {
		log.Error("AsyncCall GetEntitySrvID error: ", err, ", methodName: ", methodName)
		return err
	}

	msg := &msgdef.CallMsg{}
	msg.FromSID = e.GetIEntities().GetLocalService().GetSID()
	msg.IsSync = false
	msg.SType = uint8(sType)
	msg.SID = sid
	msg.GroupID = gID
	msg.MethodName = methodName
	msg.EntityID = e.entityID
	msg.Params = serializer.SerializeNew(args...)

	//发送给客户端
	if e.GetIEntities().GetLocalService().GetSType() == idata.ServiceGateway && sType == idata.ServiceClient {
		cli := e.GetClientSess()
		if cli != nil {
			cli.Send(msg)
			return nil
		}

		log.Error("AsyncCall failed: ClientSess is null, sType: ", sType, ", methodName: ", methodName)
		return fmt.Errorf("AsyncCall failed: ClientSess is null, sType: %d", sType)
	}

	s := iserver.GetServiceProxyMgr().GetServiceByID(sid)
	if !s.IsValid() {
		log.Error("AsyncCall service not exist, sid: ", sid, ",  sType: ", sType, " ,methodName: ", methodName)
		return fmt.Errorf("service not exist: %d", sid)
	}

	if s.IsLocal() {
		//直接发送
		is := iserver.GetLocalServiceMgr().GetLocalService(sid)
		if is == nil {
			log.Error("AsyncCall error, service is local, but not found, methodName: ", methodName)
			return fmt.Errorf("AsyncCall error, service is local, but not found")
		}

		return is.PostCallMsg(msg)
	}

	return s.SendMsg(msg)
}
