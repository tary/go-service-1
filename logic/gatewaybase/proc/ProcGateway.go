package proc

import (
	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/framework/iserver"
	"github.com/GA-TECH-SERVER/zeus/logic/gatewaybase/igateway"
	"github.com/GA-TECH-SERVER/zeus/logic/gatewaybase/userbase"

	"github.com/GA-TECH-SERVER/zeus/framework/errormsg"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
)

// PGatewayServer 是消息处理类(Processor).
// 必须实现 RegisterMsgProcFunctions(), OnClosed() 和 MsgProc*() 接口。
type PGatewayServer struct {
	sess         inet.ISession // 一般都需要包含session对象
	IServiceBase iserver.IServiceBase
	entity       iserver.IEntity
	isClosed     bool //是否已经关闭
}

// RegisterMsgProcFunctions 克隆自身并注册消息处理函数.
func (p *PGatewayServer) RegisterMsgProcFunctions(sess inet.ISession) interface{} {
	assert.True(sess != nil, "session is nil")
	result := &PGatewayServer{
		sess:         sess,
		IServiceBase: p.IServiceBase,
	}

	sess.RegMsgProc(result)
	sess.AddOnClosed(result.OnClosed)

	//log.Debugf("PGatewayServer.RegisterMsgProcFunctions, p: %p, sess: %p ", result, sess)

	return result
}

// MsgProcLoginReq  MsgProcLoginReq
func (p *PGatewayServer) MsgProcLoginReq(msg *msgdef.LoginReq) {
	log.Debugf("Begin MsgProcLoginReq, UID: %d, sess: %p", msg.UID, p.sess)

	retMsg := &msgdef.LoginResp{}

	if p.entity != nil {
		retMsg.Result = uint32(errormsg.ReturnTypeFAILRELOGIN)

		log.Error("MsgProcLoginReq:ReturnTypeFAILRELOGIN, UID: ", msg.UID)
		p.sess.Send(retMsg)
		return
	}

	ilogin, ok := p.IServiceBase.(igateway.ILoginHandler)
	if ok {
		ret := ilogin.OnLoginHandler(msg)
		if ret != msgdef.ReturnTypeSuccess {
			log.Error("MsgProc_LoginReq:OnLoginHandler err, uid: ", msg.UID, ", ret: ", ret)
			retMsg.Result = uint32(ret)
			p.sess.Send(retMsg)
			return
		}
	}

	var entity iserver.IEntity

	oldEntity := p.IServiceBase.GetEntity(msg.UID)
	if oldEntity != nil {
		log.Debugf("MsgProcLoginReq, OnReconnect, UID: %d, sess: %p, oldSess: %p", msg.UID, p.sess, oldEntity.GetClientSess())

		ireconnect, ok := oldEntity.(igateway.IReconnectHandler)
		if !ok {
			log.Errorf("MsgProcLoginReq not IReconnectHandler, UID: %d", msg.UID)

			retMsg.Result = uint32(errormsg.ReturnTypeFAILRELOGIN)
			p.sess.Send(retMsg)
			return
		}

		ret := oldEntity.PostFunctionAndWait(func() interface{} { return ireconnect.OnReconnect(p.sess) })
		retData, ok := ret.(*igateway.ReconnectData)
		if !ok {
			log.Errorf("MsgProcLoginReq OnReconnect failed, UID: %d", msg.UID)

			retMsg.Result = uint32(errormsg.ReturnTypeFAILRELOGIN)
			p.sess.Send(retMsg)
			return
		}

		if retData.Err != nil {
			log.Error("MsgProcLoginReq OnReconnect failed, UID: ", msg.UID, ", err: ", retData.Err)

			retMsg.Result = uint32(errormsg.ReturnTypeFAILRELOGIN)
			p.sess.Send(retMsg)
			return
		}

		if !retData.IsCreateEntity {
			entity = oldEntity
		}
	}

	if entity == nil {
		userBase := &userbase.UserInitData{
			Sess:    p.sess,
			Version: msg.Version,
		}

		var err error
		// 创建新玩家
		entity, err = p.IServiceBase.CreateEntityWithID("Player", msg.UID, 0, userBase, true, 0)
		if err != nil {
			log.Error("Create user failed, err: ", err, ", UID: ", msg.UID)
			retMsg.Result = uint32(errormsg.ReturnTypeFAILRELOGIN)
			p.sess.Send(retMsg)
			return
		}
	}

	//判断是否已经close
	if p.isClosed {

		iclose, ok := entity.(igateway.ICloseHandler)
		if ok {
			entity.PostFunction(func() { iclose.OnClose() })
		} else {
			log.Error("MsgProcLoginReq user not ICloseHandler, UID: ", msg.UID)
		}

		log.Errorf("MsgProcLoginReq but closed, UID: %d, p: %p, entity: %p", msg.UID, p, entity)
		return
	}

	p.entity = entity

	p.sess.SetVerified()

	//发送送登录验证成功消息
	retMsg.Result = uint32(errormsg.ReturnTypeSUCCESS)
	p.sess.Send(retMsg)

	//log.Debugf("Finish MsgProcLoginReq, UID: %d, p: %p, entity: %p", msg.UID, p, entity)
}

// OnClosed 关闭回调
func (p *PGatewayServer) OnClosed() {
	// 会话断开时动作...

	p.isClosed = true

	//log.Infof("PGatewayServer OnClosed: %d %s, p: %p, entity: %p", p.sess.GetID(), p.sess.RemoteAddr(), p, p.entity)

	if p.entity == nil {
		return
	}

	log.Debugf("PGatewayServer OnClosed, ID: %d, sess: %p", p.entity.GetEntityID(), p.sess)

	iclose, ok := p.entity.(igateway.ICloseHandler)
	if ok {
		p.entity.PostFunction(func() { iclose.OnClose() })
	} else {
		log.Error("OnClosed user not ICloseHandler, UID: ", p.entity.GetEntityID())
	}
}

// MsgProcCallMsg CallMsg消息处理
func (p *PGatewayServer) MsgProcCallMsg(msg *msgdef.CallMsg) {
	//log.Infof("MsgProcCallMsg, Seq:%d, MethodName:%s, stype: %d", msg.Seq, msg.MethodName, msg.SType)

	msg.EntityID = p.entity.GetEntityID()
	msg.IsFromClient = true

	//如果是投递本服务并且是多协程的，消息投递给实体协程
	if msg.SType == uint8(p.IServiceBase.GetSType()) && p.IServiceBase.IsMultiThread() {
		if msg.IsSync {
			retData := p.entity.PostCallMsgAndWait(msg)

			retMsg := &msgdef.CallRespMsg{}
			retMsg.Seq = msg.Seq
			retMsg.RetData = retData.Ret

			if retData.Err != nil {
				retMsg.ErrString = retData.Err.Error()
			}

			p.sess.Send(retMsg)
		} else {
			err := p.entity.PostCallMsg(msg)
			if err != nil {
				log.Error("AsyncCall err: ", err)
			}
		}
	} else if msg.SType == uint8(p.IServiceBase.GetSType()) {
		//消息投递给本服务
		msg.GroupID = p.entity.GetGroupID()
		p.postToLocalService(p.IServiceBase.GetSID(), msg)
	} else {
		//消息转发给其他服务
		msg.FromSID = p.IServiceBase.GetSID()

		srvID, gID, err := p.entity.GetEntitySrvID(msg.SType)
		if err != nil {
			log.Error("GetEntitySrvID err: ", err)
			return
		}

		//设置groupID
		msg.GroupID = gID

		proxy := iserver.GetServiceProxyMgr().GetServiceByID(srvID)
		if proxy != nil {
			if proxy.IsLocal() {
				p.postToLocalService(srvID, msg)
			} else {
				proxy.SendMsg(msg)
			}
		}
	}
}

func (p *PGatewayServer) postToLocalService(srvID uint64, msg *msgdef.CallMsg) error {
	var err error
	localS := iserver.GetLocalServiceMgr().GetLocalService(srvID)
	if localS != nil {
		if msg.IsSync {
			retData := localS.PostCallMsgAndWait(msg)

			retMsg := &msgdef.CallRespMsg{}
			retMsg.Seq = msg.Seq
			retMsg.RetData = retData.Ret

			if retData.Err != nil {
				retMsg.ErrString = retData.Err.Error()
				err = retData.Err
			}

			p.sess.Send(retMsg)
		} else {
			err = localS.PostCallMsg(msg)
			if err != nil {
				log.Error("service proxy PostCallMsg err: ", err)
			}
		}
	} else {
		//TODO:
	}

	return err
}
