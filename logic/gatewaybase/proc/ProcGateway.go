package proc

import (
	"github.com/giant-tech/go-service/base/net/inet"
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/logic/gatewaybase/igateway"

	"github.com/giant-tech/go-service/framework/errormsg"
	"github.com/giant-tech/go-service/framework/msgdef"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
)

// PGatewayServer 是消息处理类(Processor).
// 必须实现 RegisterMsgProcFunctions(), OnClosed() 和 MsgProc*() 接口。
type PGatewayServer struct {
	sess         inet.ISession        // 一般都需要包含session对象
	IServiceBase iserver.IServiceBase // 所属的服务
	entity       iserver.IEntity      // 实体
	group        iserver.IEntityGroup // 实体所属的group，如果不属于任何group，则为nil

	proxyEntity iserver.IEntity
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
	if !ok {
		retMsg.Result = uint32(errormsg.ReturnTypeFAILRELOGIN)

		log.Error("MsgProcLoginReq:ReturnTypeFAILRELOGIN, UID: ", msg.UID)
		p.sess.Send(retMsg)
		return
	}

	retData := ilogin.OnLoginHandler(p.sess, msg)
	if retData.Msg.Result != uint32(msgdef.ReturnTypeSuccess) {
		p.sess.Send(retData.Msg)
		return
	}

	p.entity = retData.Entity
	p.group = retData.Group

	if p.group != nil {
		p.proxyEntity = p.group
	} else {
		p.proxyEntity = p.entity
	}

	p.sess.SetVerified()

	//发送送登录验证成功消息
	retMsg.Result = uint32(errormsg.ReturnTypeSUCCESS)
	p.sess.Send(retData.Msg)

	//log.Debugf("Finish MsgProcLoginReq, UID: %d, p: %p, entity: %p", msg.UID, p, entity)
}

// OnClosed 关闭回调
func (p *PGatewayServer) OnClosed() {
	// 会话断开时动作...

	//log.Infof("PGatewayServer OnClosed: %d %s, p: %p, entity: %p", p.sess.GetID(), p.sess.RemoteAddr(), p, p.entity)

	if p.proxyEntity == nil {
		return
	}

	log.Debugf("PGatewayServer OnClosed, ID: %d, sess: %p", p.entity.GetEntityID(), p.sess)

	iclose, ok := p.entity.(igateway.ICloseHandler)
	if ok {
		p.proxyEntity.PostFunction(func() { iclose.OnClose() })
	} else {
		log.Error("OnClosed user not ICloseHandler, UID: ", p.entity.GetEntityID())
	}

	p.entity = nil
	p.group = nil
	p.proxyEntity = nil
}

// MsgProcCallMsg CallMsg消息处理
func (p *PGatewayServer) MsgProcCallMsg(msg *msgdef.CallMsg) {
	//log.Infof("MsgProcCallMsg, Seq:%d, MethodName:%s, stype: %d", msg.Seq, msg.MethodName, msg.SType)

	msg.EntityID = p.entity.GetEntityID()
	msg.IsFromClient = true

	//如果是本服务，则把消息投递给实体或者服务
	if msg.SType == uint8(p.IServiceBase.GetSType()) {
		// 如果为多协程，则投递给对应的实体
		if p.IServiceBase.IsMultiThread() {
			if msg.IsSync {
				retData := p.proxyEntity.PostCallMsgAndWait(msg)

				retMsg := &msgdef.CallRespMsg{}
				retMsg.Seq = msg.Seq
				retMsg.RetData = retData.Ret

				if retData.Err != nil {
					retMsg.ErrString = retData.Err.Error()
				}

				p.sess.Send(retMsg)
			} else {
				err := p.proxyEntity.PostCallMsg(msg)
				if err != nil {
					log.Error("AsyncCall err: ", err)
				}
			}
		} else {
			//消息投递给本服务
			p.postToLocalService(p.IServiceBase.GetSID(), msg)
		}
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
