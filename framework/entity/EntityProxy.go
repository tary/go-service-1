package entity

import (
	"fmt"

	"github.com/GA-TECH-SERVER/zeus/base/serializer"
	"github.com/GA-TECH-SERVER/zeus/framework/idata"
	"github.com/GA-TECH-SERVER/zeus/framework/iserver"
	redis "github.com/GA-TECH-SERVER/zeus/framework/logicredis"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef" //sess "github.com/GA-TECH-SERVER/zeus/netentity"

	"github.com/cihub/seelog"
)

// EProxy 实体的一个代理，可以方便的传递
type EProxy struct {
	EntityID uint64
}

// NewEntityProxy 创建一个新的实体代理
func NewEntityProxy(entityID uint64) *EProxy {
	return &EProxy{
		EntityID: entityID,
	}
}

// SyncCall 同步调用，等待返回
func (e *EProxy) SyncCall(stype idata.ServiceType, retPtr interface{}, methodName string, args ...interface{}) error {
	sid, gid, err := redis.GetEntitySrvUtil(e.EntityID).GetSrvInfo(uint8(stype))
	if err != nil {
		return err
	}

	s := iserver.GetServiceProxyMgr().GetServiceByID(sid)
	if !s.IsValid() {
		return fmt.Errorf("service not exist: %d", sid)
	}

	msg := &msgdef.CallMsg{}

	msg.SType = uint8(stype)
	msg.SID = sid
	msg.GroupID = gid
	msg.MethodName = methodName
	msg.IsSync = true
	msg.EntityID = e.EntityID
	msg.Params = serializer.SerializeNew(args...)

	if s.IsLocal() {
		//直接发送
		is := iserver.GetLocalServiceMgr().GetLocalService(sid)
		if is == nil {
			seelog.Error("")
			return fmt.Errorf("")
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
func (e *EProxy) AsyncCall(stype idata.ServiceType, methodName string, args ...interface{}) error {
	sid, gid, err := redis.GetEntitySrvUtil(e.EntityID).GetSrvInfo(uint8(stype))
	if err != nil {
		return err
	}

	s := iserver.GetServiceProxyMgr().GetServiceByID(sid)
	if s.IsValid() {
		return fmt.Errorf("service not exist: %d", sid)
	}

	msg := &msgdef.CallMsg{}

	msg.IsSync = false
	msg.SID = sid
	msg.GroupID = gid
	msg.EntityID = e.EntityID
	msg.MethodName = methodName
	msg.Params = serializer.SerializeNew(args...)

	if s.IsLocal() {
		//直接发送
		is := iserver.GetLocalServiceMgr().GetLocalService(sid)
		if is == nil {
			seelog.Error("")
			return fmt.Errorf("")
		}

		return is.PostCallMsg(msg)
	}
	s.SendMsg(msg)

	return nil
}
