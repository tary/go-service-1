package baseproc

import (
	"github.com/GA-TECH-SERVER/zeus/base/net/inet"

	assert "github.com/aurelien-rainone/assertgo"
	"github.com/cihub/seelog"
)

// BProcServer 基本处理server
type BProcServer struct {
	sess       inet.ISession
	serverID   uint64
	serverType int32
}

// RegisterMsgProcFunctions 克隆自身并注册消息处理函数.
func (p *BProcServer) RegisterMsgProcFunctions(sess inet.ISession) interface{} {
	assert.True(sess != nil, "session is nil")

	seelog.Debug("BProcServer, RegisterMsgProcFunctions")

	result := &BProcServer{
		sess: sess,
	}

	sess.RegMsgProc(result)
	sess.AddOnClosed(result.OnClosed)

	return result
}

// OnClosed 关闭回调
func (p *BProcServer) OnClosed() {
	// 会话断开时动作...
	/*iserver.GetSrvInst().GetServerCtrl().OnDisconnect(w.serverID, (uint8)(w.serverType))

	sdsess.DeleteSession(w.serverID)
	seelog.Info("BProcServer closed:", w.serverID, " ", w.serverType)
	*/
}
