package server

import (
	"github.com/giant-tech/go-service/framework/net/inet"
	connhandler "github.com/giant-tech/go-service/framework/net/server/internal/conn_handler"
	"github.com/giant-tech/go-service/framework/net/server/internal/listener"
	msgprocset "github.com/giant-tech/go-service/framework/net/server/internal/msg_proc_set"

	"golang.org/x/time/rate"
)

// Server 结构体
type Server struct {
	listener    listener.IListener
	connHandler *connhandler.ConnHandler
}

// New 创建一个服务服.
// 自动开始监听。
// protocol 支持："kcp", "tcp", "tcp+kcp".
// "tcp+kcp"可接受kcp客户端，也可接受tcp客户端。
// addr 形如：":80", "1.2.3.4:80"
// maxConns 是最大连接数
func New(protocol string, addr string, maxConns int) (*Server, error) {
	srv := &Server{
		connHandler: connhandler.New(),
	}

	//srv.AddMsgProc(&baseproc.BaseProcServer{})

	var err error
	// 接受新连接时会 go sessMgr.HandleConn(), 运行 session.Start().
	srv.listener, err = listener.NewListener(protocol, addr, maxConns)
	return srv, err
}

// Run 运行服务器.
func (s *Server) Run() {
	if !s.connHandler.HasMsgProc() {
		// 防止直接创建, 忘记 AddMsgProc().
		// 如果是生成的代码，会自动 AddMsgProc().
		panic("No MsgProc!")
	}

	s.listener.Run(s.connHandler)
}

// Close 结束监听, 结束所有消息处理.
func (s *Server) Close() {
	s.listener.Close()
}

// SetSessEvtSink 设置一个会话事件接收器.
// Depricated. 应该使用 MsgProc 的 OnClosed(), 更简单。
func (s *Server) SetSessEvtSink(sink inet.ISessEvtSink) {
	s.connHandler.SetSessEvtSink(sink)
}

// AddMsgProc 添加消息处理
func (s *Server) AddMsgProc(msgProc msgprocset.IMsgProc) {
	s.connHandler.AddMsgProc(msgProc)
}

//SetEncrypt 设置是否加密
func (s *Server) SetEncrypt(isEncrypt bool) {
	s.connHandler.SetEncrypt(isEncrypt)
}

//GetListenPort 获取监听端口
func (s *Server) GetListenPort() string {
	return s.listener.GetPort()
}

// SetVerifyMsgID 设置会话的验证消息ID.
// 强制会话必须验证，会话的第1个消息将做为验证消息，必须是指定消息号。
// 应用的MsgProc处理器必须调用 session.SetVerified(), 不然连接将被强制关闭。
func (s *Server) SetVerifyMsgID(verifyMsgID inet.MsgID) {
	s.connHandler.SetVerifyMsgID(verifyMsgID)
}

// SetBytePerSecLimiter 设置每秒接收字节数限制.
// r(rate) 为每秒字节数。
// b(burst) 为峰值字节数。
func (s *Server) SetBytePerSecLimiter(r rate.Limit, b int) {
	// connHandler 会在每个Session创建时设置限流
	s.connHandler.SetBytePerSecLimiter(r, b)
}

// SetQueryPerSecLimiter 设置每秒接收请求数限制.
// r(rate) 为每秒请求数。
// b(burst) 为峰值请求数。
// 必须在 Start() 之前设置，避免 DataRace.
func (s *Server) SetQueryPerSecLimiter(r rate.Limit, b int) {
	// connHandler 会在每个Session创建时设置限流
	s.connHandler.SetQueryPerSecLimiter(r, b)
}
