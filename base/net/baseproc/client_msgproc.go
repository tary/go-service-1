package baseproc

import (
	"github.com/GA-TECH-SERVER/zeus/base/net/baseproc/basemsg"
	"github.com/GA-TECH-SERVER/zeus/base/net/client"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"
	"github.com/GA-TECH-SERVER/zeus/framework/sd/sdsess"

	assert "github.com/aurelien-rainone/assertgo"
	"github.com/cihub/seelog"
)

// BPClient 基本的消息处理
type BPClient struct {
	sess       *client.Session
	serverType int32
	serverID   uint64
}

// RegProcClientBaseProc 注册处理函数
func RegProcClientBaseProc(sess *client.Session) {
	assert.True(sess != nil, "session is nil")
	w := &BPClient{
		sess: sess,
	}

	// [ServerToClient] 注册接收的消息。需要从ID创建消息。
	sess.RegMsgProc(w)

	sess.AddOnClosed(w.OnClosed)
}

// OnClosed 关闭回调
func (w *BPClient) OnClosed() {
	seelog.Info("BPClient.OnClosed")

	//iserver.GetSrvInst().GetServerCtrl().OnDisconnect(w.serverID, (uint8)(w.serverType))

	sdsess.DeleteSession(w.serverID)
}

// MsgProcClientVerifyResp 消息请求回应
func (w *BPClient) MsgProcClientVerifyResp(resp *msgdef.ClientVerifyResp) {
	if resp.Result != 0 {
		//验证失败
		seelog.Error("MsgProcClientVerifyResp verify failed:", resp)
		w.sess.Close()
		return
	}
	w.serverType, w.serverID = int32(resp.ServerType), resp.ServerID
	seelog.Info("BPClient.MsgProcClientVerifyResp: ", w.serverType, " ", w.serverID)

	sdsess.AddSession(w.serverID, w.serverType, w.sess)

	verifyData := &basemsg.VerifyData{}
	verifyData.ServerType = resp.ServerType
	verifyData.ServerID = resp.ServerID

	w.sess.Emit("verified", verifyData)

	//iserver.GetSrvInst().GetServerCtrl().OnConnect(w.serverID, (uint8)(w.serverType))
}
