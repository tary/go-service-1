package connhandler

import (
	"net"

	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	msgprocset "github.com/GA-TECH-SERVER/zeus/base/net/server/internal/msg_proc_set"
	"github.com/GA-TECH-SERVER/zeus/base/net/server/internal/session"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"

	"golang.org/x/time/rate"
)

// IMsgCreatorGetor 消息创造
type IMsgCreatorGetor interface {
	GetMsgCreator(uint32) inet.IMsgCreator
}

// ConnHandler 连接处理
type ConnHandler struct {
	msgProcSet  *msgprocset.MsgProcSet
	sessEvtSink inet.ISessEvtSink

	sessVerifyMsgID inet.MsgID // 验证消息ID
	isEncrypt       bool       //是否需要加密
	// 限流
	bpsLimiter *rate.Limiter // 限制每秒接收字节数
	qpsLimiter *rate.Limiter // 限制每秒接收请求数
}

// New 新建
func New() *ConnHandler {
	return &ConnHandler{
		sessVerifyMsgID: inet.MsgID(msgdef.ClientVerifyReqMsgID),
		msgProcSet:      msgprocset.New(),
	}
}

// HandleConn 处理连接.
// 创建 Session, 并且开始会话协程。
func (h *ConnHandler) HandleConn(conn net.Conn) {
	sess := session.New(conn, h.isEncrypt, h.sessEvtSink, false)

	// 在会话上注册所有处理函数
	if h.msgProcSet != nil {
		h.msgProcSet.RegisterToSession(sess)
	}

	sess.On("verified", func(v interface{}) {
		sess.Emit("initAfterVerified", v)
	})

	sess.SetVerifyMsgID(h.sessVerifyMsgID)

	// 设置限流
	h.setLimitersToSession(sess)

	sess.Start()
}

// SetSessEvtSink 设置sess
func (h *ConnHandler) SetSessEvtSink(sink inet.ISessEvtSink) {
	h.sessEvtSink = sink
}

//SetEncrypt 设置是否加密，默认不加密
func (h *ConnHandler) SetEncrypt(isEncrypt bool) {
	h.isEncrypt = isEncrypt
}

// AddMsgProc 添加消息处理
func (h *ConnHandler) AddMsgProc(msgProc msgprocset.IMsgProc) {
	h.msgProcSet.AddMsgProc(msgProc)
}

// HasMsgProc 是否有消息处理
func (h *ConnHandler) HasMsgProc() bool {
	return !h.msgProcSet.IsEmpty()
}

// SetVerifyMsgID 设置会话的验证消息.
// 强制会话必须验证，会话的第1个消息将做为验证消息，消息类型必须为输入类型.
func (h *ConnHandler) SetVerifyMsgID(verifyMsgID inet.MsgID) {
	h.sessVerifyMsgID = verifyMsgID
}

// SetBytePerSecLimiter 设置每秒接收字节数限制.
// r(rate) 为每秒字节数。
// b(burst) 为峰值字节数。
func (h *ConnHandler) SetBytePerSecLimiter(r rate.Limit, b int) {
	h.bpsLimiter = rate.NewLimiter(r, b)
}

// SetQueryPerSecLimiter 设置每秒接收请求数限制.
// r(rate) 为每秒请求数。
// b(burst) 为峰值请求数。
// 必须在 Start() 之前设置，避免 DataRace.
func (h *ConnHandler) SetQueryPerSecLimiter(r rate.Limit, b int) {
	h.qpsLimiter = rate.NewLimiter(r, b)
}

// setLimitersToSession 设置session限流
func (h *ConnHandler) setLimitersToSession(sess *session.Session) {
	bps := h.bpsLimiter
	if bps != nil {
		sess.SetBytePerSecLimiter(bps.Limit(), bps.Burst())
	}

	qps := h.qpsLimiter
	if qps != nil {
		sess.SetQueryPerSecLimiter(qps.Limit(), qps.Burst())
	}
}
