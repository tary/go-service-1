package client

import (
	"net"

	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/net/inet"
	"github.com/giant-tech/go-service/framework/net/internal"

	"golang.org/x/time/rate"
)

// Session 包装内部会话，并提供额外的客户端功能
type Session struct {
	sess inet.ISession
}

// NewSession 新建一个session
func NewSession(conn net.Conn, isIdip bool) *Session {
	sess := &Session{
		sess: internal.NewSession(conn, false, true, isIdip),
	}

	sess.sess.SetOnClosed(sess.onClosed)

	return sess
}

// onClosed 关闭回调
func (s *Session) onClosed() {

}

// Send 发送
func (s *Session) Send(msg imsg.IMsg) error {
	return s.sess.Send(msg)
}

// SendRaw 发送原始数据
func (s *Session) SendRaw(buff []byte) error {
	return s.sess.SendRaw(buff)
}

// SetEncrypt 设置加密
func (s *Session) SetEncrypt(isEncrypt bool) {
	s.sess.SetEncrypt(isEncrypt)
}

// EncodeMsg 编码消息
func (s *Session) EncodeMsg(msg imsg.IMsg) ([]byte, error) {
	return s.sess.EncodeMsg(msg)
}

// CompressAndEncrypt 压缩并加密数据
func (s *Session) CompressAndEncrypt(buf []byte) ([]byte, error) {
	return s.sess.CompressAndEncrypt(buf)
}

// Start session开始
func (s *Session) Start() {
	s.sess.Start()
}

// Close session关闭
func (s *Session) Close() {
	s.sess.Close()
}

// IsClosed session 是否关闭
func (s *Session) IsClosed() bool {
	return s.sess.IsClosed()
}

// RegMsgProcFunc 注册消息处理函数.
// 3个参数为：消息ID, 消息创建函数，消息处理函数。
// 必须在Start()之前。
func (s *Session) RegMsgProcFunc(msgID inet.MsgID, procFunc func(imsg.IMsg)) {
	s.sess.RegMsgProcFunc(msgID, procFunc)
}

// RegMsgProc 注册类中所有消息处理函数
func (s *Session) RegMsgProc(proc interface{}) {
	s.sess.RegMsgProc(proc)
}

// RegIdipMsgProcFunc 注册idip消息处理函数.
// 3个参数为：消息ID, 消息创建函数，消息处理函数。
// 必须在Start()之前。
func (s *Session) RegIdipMsgProcFunc(msgID inet.IdipMsgID, procFunc func(interface{})) {
	// 注册接收的消息。需要从ID创建消息。
	s.sess.RegIdipMsgProcFunc(msgID, procFunc)
}

// SendIdip idip发送
func (s *Session) SendIdip(msg interface{}) {
	s.sess.SendIdip(msg)
}

// On 开始运行
func (s *Session) On(evt string, f func(interface{})) {
	s.sess.On(evt, f)
}

// Emit 发射
func (s *Session) Emit(evt string, p interface{}) {
	s.sess.Emit(evt, p)
}

// AddOnClosed 添加onclosed回调
func (s *Session) AddOnClosed(f func()) {
	s.sess.AddOnClosed(f)
}

// GetID 获取sess id
func (s *Session) GetID() uint64 {
	return s.sess.GetID()
}

// RemoteAddr sess远程地址
func (s *Session) RemoteAddr() string {
	return s.sess.RemoteAddr()
}

// ResetHb 重置心跳
func (s *Session) ResetHb() {
	panic("only server can call")
}

// SetOnClosed 设置关闭回调
func (s *Session) SetOnClosed(func()) {
	panic("only server can call")
}

// SetVerifyMsgID 设置
func (s *Session) SetVerifyMsgID(inet.MsgID) {
	panic("only server can call")
}

// SetVerified 设置验证过
func (s *Session) SetVerified() {
	panic("only server can call")
}

// SetBytePerSecLimiter set
func (s *Session) SetBytePerSecLimiter(r rate.Limit, b int) {
	panic("only server can call")
}

// SetQueryPerSecLimiter set
func (s *Session) SetQueryPerSecLimiter(r rate.Limit, b int) {
	panic("only server can call")
}
