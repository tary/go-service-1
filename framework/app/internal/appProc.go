package internal

import (
	"fmt"

	"github.com/giant-tech/go-service/framework/errormsg"
	"github.com/giant-tech/go-service/framework/idata"
	"github.com/giant-tech/go-service/framework/iserver"
	"github.com/giant-tech/go-service/framework/msgdef"
	"github.com/giant-tech/go-service/framework/net/baseproc/basemsg"
	"github.com/giant-tech/go-service/framework/net/inet"
	"github.com/giant-tech/go-service/framework/service"

	assert "github.com/aurelien-rainone/assertgo"
	"github.com/cihub/seelog"
)

// ProcApp 是消息处理类(Processor).
// 必须实现 MsgProc*() 接口。
type ProcApp struct {
	Sess        inet.ISession // 一般都需要包含session对象
	remoteAppID uint64
}

// RegisterMsgProcFunctions 注册消息处理
func (p *ProcApp) RegisterMsgProcFunctions(sess inet.ISession) interface{} {
	assert.True(sess != nil, "session is nil")
	proc := &ProcApp{
		Sess: sess,
	}

	sess.RegMsgProc(proc)

	sess.AddOnClosed(proc.OnClosed)

	return proc
}

// OnClosed app关闭回调
func (p *ProcApp) OnClosed() {
	// 会话断开时动作...
	seelog.Infof("OnClosed sessID:%d, remote addr:%s", p.Sess.GetID(), p.Sess.RemoteAddr())

	//通知app，远程服务不可用了
	service.GetServiceProxyMgr().OnServiceClosed(p.remoteAppID)

	//删除远程服务的session
	service.GetServiceProxyMgr().DelServiceByAppID(p.remoteAppID)

}

// MsgProcClientVerifyReq 客户端消息请求
func (p *ProcApp) MsgProcClientVerifyReq(msg *msgdef.ClientVerifyReq) {
	seelog.Debug("MsgProcClientVerifyReq, token: ", msg.Token)

	resp := &msgdef.ClientVerifyResp{
		ServerID:   MyApp.appID,
		ServerType: 1, //temp
	}

	//期望其它服务器发来的token（serverid）和本服是一致的，否则验证失败
	//serverID, _ := strconv.ParseUint(msg.Token, 10, 64)

	if msg.Token != fmt.Sprintf("%d", resp.ServerID) {
		resp.Result = uint32(errormsg.ReturnTypeTOKENINVALID)
		p.Sess.Send(resp)
		return
	}

	p.remoteAppID = msg.ServerID

	p.Sess.SetVerified()

	//ClientVerifyResp消息处理
	resp.Result = uint32(errormsg.ReturnTypeSUCCESS)
	p.Sess.Send(resp)

	iserver.GetApp().GetAppNet().InsertSrvSess(msg.ServerID, p.Sess)
}

// MsgProcClientVerifyResp 客户端消息回应
func (p *ProcApp) MsgProcClientVerifyResp(resp *msgdef.ClientVerifyResp) {
	if resp.Result != uint32(errormsg.ReturnTypeSUCCESS) {
		//验证失败
		seelog.Error("MsgProcClientVerifyResp verify failed:", resp)
		p.Sess.Close()
		return
	}
	//p.serverType, p.serverID = int32(resp.ServerType), resp.ServerID
	seelog.Info("BPClient.MsgProcClientVerifyResp: " /*, w.serverType, " ", w.serverID*/)

	//sdsess.AddSession(w.serverID, w.serverType, w.sess)

	verifyData := &basemsg.VerifyData{}
	verifyData.ServerType = resp.ServerType
	verifyData.ServerID = resp.ServerID

	p.Sess.Emit("verified", verifyData)

	//记录连上的远程appID,放在这里是否合适？
	p.remoteAppID = resp.ServerID
	iserver.GetApp().GetAppNet().OnServerConnected(resp.ServerID)

}

// MsgProcCallMsg msg
func (p *ProcApp) MsgProcCallMsg(msg *msgdef.CallMsg) {
	//seelog.Infof("MsgProcCallMsg, ToServiceID:%d,Seq:%d, MethodName:%s",
	//	msg.ToServiceID, msg.Seq, msg.MethodName)

	s := service.GetLocalServiceMgr().GetLocalService(msg.SID)
	if s == nil {
		return
	}

	if msg.IsSync {
		retData := s.PostCallMsgAndWait(msg)

		retMsg := &msgdef.CallRespMsg{}
		retMsg.Seq = msg.Seq
		retMsg.RetData = retData.Ret

		if retData.Err != nil {
			retMsg.ErrString = retData.Err.Error()
		}

		if msg.IsFromClient {
			var err error
			fmsg := &msgdef.ForwardToClientMsg{}
			fmsg.ServiceID = msg.FromSID
			fmsg.EntityID = msg.EntityID
			fmsg.MsgData, err = p.Sess.EncodeMsg(retMsg)
			if err != nil {
				seelog.Error("EncodeMsg error: ", err)
				return
			}

			p.Sess.Send(fmsg)
		} else {
			p.Sess.Send(retMsg)
		}
	} else {
		err := s.PostCallMsg(msg)
		if err != nil {
			seelog.Error("AsyncCall err: ", err)
		}
	}
}

// MsgProcCallRespMsg resp
func (p *ProcApp) MsgProcCallRespMsg(msg *msgdef.CallRespMsg) {
	seelog.Infof("MsgProcCallRespMsg ErrString:%s, Seq:%d",
		msg.ErrString, msg.Seq)

	call := MyApp.GetPendingCall(msg.Seq)
	if call == nil {
		seelog.Error("MsgProcCallRespMsg, can't find pending call, seq: ", msg.Seq)
		return
	}

	MyApp.DelPendingCall(msg.Seq)

	retData := &idata.RetData{}
	if len(msg.ErrString) > 0 {
		retData.Err = fmt.Errorf("%s", msg.ErrString)
	} else {
		retData.Ret = msg.RetData
	}

	call.RetChan <- retData
}

// MsgProcForwardToClientMsg 转发
func (p *ProcApp) MsgProcForwardToClientMsg(msg *msgdef.ForwardToClientMsg) {

	is := iserver.GetLocalServiceMgr().GetLocalService(msg.ServiceID)
	if is == nil {
		seelog.Error("service not found, serviceID: ", msg.ServiceID)
		return
	}

	e := is.GetEntity(msg.EntityID)
	if e == nil || e.GetClientSess() == nil {
		seelog.Error("entity not found or GetClientSess is nil")
		return
	}

	e.GetClientSess().SendRaw(msg.MsgData)
}
