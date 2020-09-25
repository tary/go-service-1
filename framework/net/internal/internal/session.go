package internal

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/msgdef"
	"github.com/giant-tech/go-service/framework/net/inet"
	"github.com/giant-tech/go-service/framework/net/internal/internal/msgenc"
	"github.com/giant-tech/go-service/framework/net/internal/internal/msghdl"
	"github.com/giant-tech/go-service/framework/net/internal/internal/sflist"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"golang.org/x/time/rate"
)

// sessionIDGen 产生sesssionID
var sessionIDGen atomic.Uint64

// NewSession 新建session
func NewSession(conn net.Conn, encryEnabled bool, isClient bool, isIdip bool) *Session {

	sess := &Session{
		id: sessionIDGen.Inc(),

		conn:    conn,
		sendBuf: sflist.NewSafeList(),

		hbTimerInterval: time.Duration(viper.GetInt("Config.HBTimer")) * time.Second,
		hbEnabled:       viper.GetBool("Config.HeartBeat"),
		isClient:        isClient,

		_encryEnabled: encryEnabled,

		bpsLimiter: rate.NewLimiter(rate.Inf, 0),
		qpsLimiter: rate.NewLimiter(rate.Inf, 0),

		isIdip: isIdip,
		evtMap: make(map[string]func(interface{})),
	}

	if isIdip {
		sess.msgHdlr = msghdl.NewIdipHandler()
	} else {
		sess.msgHdlr = msghdl.New()
	}

	sess.hbTimer = time.NewTimer(sess.hbTimerInterval)
	sess.ctx, sess.ctxCancel = context.WithCancel(context.Background())
	return sess
}

// TODO: thread-safe

// Session 代表一个网络连接
type Session struct {
	id      uint64           //唯一id
	conn    net.Conn         // 底层连接对象
	sendBuf *sflist.SafeList // 缓存待发送的数据

	ctx       context.Context
	ctxCancel context.CancelFunc

	msgHdlr msghdl.IMessgeHandler

	hbTimer         *time.Timer
	hbTimerInterval time.Duration
	hbEnabled       bool
	isClient        bool

	_isClosed     atomic.Bool
	_isError      atomic.Bool
	_encryEnabled bool

	onClosedFuncs []func() // 关闭时依次调用

	// 关闭事件处理器
	onClosed func()

	// 会话验证功能，仅服务器有用。可以设置会话需要验证，第1个消息必须为验证消息。
	verifyMsgID inet.MsgID  // 验证的消息id
	needVerify  bool        // 是否需要验证
	isVerified  atomic.Bool // 是否已通过验证。由消息处理器设置。
	evtMap      map[string]func(interface{})

	// idip消息结构不一样
	isIdip bool //是不是idip

	// 限流
	bpsLimiter *rate.Limiter // 限制每秒接收字节数
	qpsLimiter *rate.Limiter // 限制每秒接收请求数
}

// Start 验证完成
func (sess *Session) Start() {
	if !sess.isIdip {
		sess.regHeartBeat()
	}

	if sess.hbEnabled && !sess.isIdip {
		sess.enableHeartBeat()
	}

	go sess.recvLoop() // 协程中会调用消息处理器
	go sess.sendLoop()
}

// SetEncrypt 设置加密
func (sess *Session) SetEncrypt(isEncrypt bool) {
	sess._encryEnabled = isEncrypt
}

// Send 发送消息，立即返回.
func (sess *Session) Send(msg imsg.IMsg) error {
	if msg == nil {
		return fmt.Errorf("msg is nil")
	}

	if sess.IsClosed() {
		log.Warnf("Send after sess close localAddr:%s remoteAddr:%s %s %s",
			sess.conn.LocalAddr(), sess.conn.RemoteAddr(), reflect.TypeOf(msg),
			fmt.Sprintf("%s", msg)) // log是异步的，所以 msg 必须复制下。
		return fmt.Errorf("session closed")
	}

	// 队列已满, 表示客户端处理太慢，断开。
	// XXX
	//	if len(sess.sendBufC) >= cap(sess.sendBufC) {
	//		log.Error("Close slow Session: ", sess.conn.RemoteAddr())
	//		sess.Close()
	//		return
	//	}

	msgBuf, err := sess.EncodeMsg(msg)
	if err != nil {
		log.Error("Encode message error in Send(): ", err)
		return err
	}

	sess.sendBuf.Put(msgBuf)

	return nil
}

// SendRaw 发送数据，立即返回.
func (sess *Session) SendRaw(buff []byte) error {
	if sess.IsClosed() {
		log.Warnf("Send after sess close remoteAddr: %s, localAddr: %s", sess.conn.RemoteAddr(), sess.conn.LocalAddr())
		return fmt.Errorf("session closed")
	}

	sess.sendBuf.Put(buff)

	return nil
}

// EncodeMsg 编码信息
func (sess *Session) EncodeMsg(msg imsg.IMsg) ([]byte, error) {
	msgID, err := msgdef.GetMsgDef().GetMsgIDByType(msg)
	if err != nil {
		// 应该用 msg2id.RegMsg2ID()注册才行
		log.Errorf("message '%s' is not registered", reflect.TypeOf(msg))
		return nil, fmt.Errorf("message '%s' is not registered", reflect.TypeOf(msg))
	}

	msgBuf, err := msgenc.EncodeMsg(msg, msgID)
	if err != nil {
		log.Error("Encode message error in EncodeMsg(): ", err)
		return nil, err
	}

	return msgBuf, nil
}

// CompressAndEncrypt 压缩和加密已序列化消息.
// 输入消息带头部长度和消息ID。
func (sess *Session) CompressAndEncrypt(buf []byte) ([]byte, error) {
	return msgenc.CompressAndEncrypt(buf, false, sess._encryEnabled, sess.isClient)
}

// ResetHb  记录心跳状态
func (sess *Session) ResetHb() {
	sess.hbTimer.Reset(sess.hbTimerInterval)
}

// hbLoopCheck hbLoopCheck
func (sess *Session) hbLoopCheck() {

	defer func() {
		if err := recover(); err != nil {
			log.Error("hbLoopCheck panic:", err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	for {
		select {
		case <-sess.ctx.Done():
			sess.hbTimer.Stop()
			return
		case <-sess.hbTimer.C:
			log.Error("sess heart tick expired ", sess.conn.RemoteAddr())

			sess._isError.Store(true)
			sess.hbTimer.Stop()
			sess.Close()
			return
		}
	}
}

// hbLoopSend hbLoopSend
func (sess *Session) hbLoopSend() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("hbLoopCheck panic:", err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	hbTickerInterval := time.Duration(viper.GetInt("Config.HBTicker")) * time.Second
	hbTicker := time.NewTicker(hbTickerInterval)
	for {
		select {
		case <-hbTicker.C:
			if sess.IsClosed() {
				hbTicker.Stop()
				return
			}
			sess.Send(&msgdef.Ping{})
		}
	}
}

// enableHeartBeat enableHeartBeat
func (sess *Session) enableHeartBeat() {
	sess.ResetHb()
	go sess.hbLoopCheck()

	if sess.isClient {
		go sess.hbLoopSend()
	}
}

// regHeartBeat regHeartBeat
func (sess *Session) regHeartBeat() {
	p := NewHeartBeatProc(sess)
	sess.RegMsgProc(p)
}

// recvLoop recvLoop
func (sess *Session) recvLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("recvLoop panic:", err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	var err error
	assert.True(err == nil, "init error is not nil")

	for {
		select {
		case <-sess.ctx.Done():
			sess.callonClosed()
			return
		default:
		} // select

		// readAndHandleOneMsg 有限流
		if err = sess.readAndHandleOneMsg(); err != nil {
			break
		}
	} // for

	assert.True(err != nil, "must be error")

	if !sess.needVerify {
		sess._isError.Store(true)
	}

	// 底层检测到连接断开，可能是客户端主动Close或客户端断网超时
	log.Errorf("recvLoop error error: %s, remoteAddr: %s,  localAddr: %s, sessionid: %d", err, sess.conn.RemoteAddr(), sess.conn.LocalAddr(), sess.id)
	sess.Close()

	sess.callonClosed()

	return
}

// sendLoop sendLoop
func (sess *Session) sendLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("sendLoop panic:", err, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(err)
			}
		}
	}()

	for {
		select {
		case <-sess.ctx.Done():
			return
		case <-sess.sendBuf.HasDataC: // 可能有数据了
			if sess.sendAllBufferedData() {
				continue // 正常，继续
			}
			assert.True(sess.IsClosed())
			return // 出错了，退出
		}
	}
}

// sendAllBufferedData 发送所有缓存数据，成功则返回true, 失败则关闭会话并返回false.
func (sess *Session) sendAllBufferedData() bool {
	for {
		data, err := sess.sendBuf.Pop()
		if err != nil {
			return true // 已取完
		}

		if !sess.sendData(data.([]byte)) {
			assert.True(sess.IsClosed())
			return false // 出错退出
		}
	}
}

// sendData 发送数据，成功返回true, 失败则关闭会话并返回false.
func (sess *Session) sendData(buf []byte) bool {
	var msgBuf []byte
	var err error

	if sess.isIdip {
		msgBuf, err = msgenc.CompressAndEncryptIdip(buf, true, sess._encryEnabled)
	} else {
		msgBuf, err = msgenc.CompressAndEncrypt(buf, false, sess._encryEnabled, sess.isClient)
	}

	if err != nil {
		log.Error("compress and encrypt message error: ", err)
		return true // Todo: 是否应该出错关闭？
	}

	_, err = sess.conn.Write(msgBuf)
	if err == nil {
		return true
	}

	sess._isError.Store(true)
	if sess.IsClosed() {
		return false
	}

	log.Error("send message error ", err)
	sess.Close()
	return false
}

// Close 关闭.
// 所有发送完成后才关闭。或2s后强制关闭。
func (sess *Session) Close() {
	if !sess._isClosed.CAS(false, true) {
		return
	}

	sess.hbTimer.Stop()

	if !sess._isError.Load() {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Error("Close _isError panic:", err, string(debug.Stack()))
					if viper.GetString("Config.Recover") == "0" {
						panic(err)
					}
				}
			}()

			closeTicker := time.NewTicker(100 * time.Millisecond)
			defer closeTicker.Stop()
			closeTimer := time.NewTimer(2 * time.Second)
			defer closeTimer.Stop()
			for {
				select {
				case <-closeTimer.C:
					sess.ctxCancel()
					sess.conn.Close()
					return
				case <-closeTicker.C:
					if sess.sendBuf.IsEmpty() {
						sess.ctxCancel()
						sess.conn.Close()
						return
					}
				}
			}
		}()
	} else {
		sess.ctxCancel()
		sess.conn.Close()
	}
}

// callonClosed 由接收消息线程调用onClosed
func (sess *Session) callonClosed() {
	//log.Debug("callonClosed")

	if sess.onClosed != nil {
		for _, f := range sess.onClosedFuncs {
			f()
		}

		sess.onClosed()
	}
}

// RemoteAddr 远程地址
func (sess *Session) RemoteAddr() string {
	return sess.conn.RemoteAddr().String()
}

// IsClosed 返回sess是否已经关闭
func (sess *Session) IsClosed() bool {
	return sess._isClosed.Load()
}

// RegMsgProcFunc 注册消息处理函数.
func (sess *Session) RegMsgProcFunc(msgID inet.MsgID, procFun func(imsg.IMsg)) {
	sess.msgHdlr.RegMsgProcFunc(msgID, procFun)
}

// RegMsgProc 注册消息处理函数.
func (sess *Session) RegMsgProc(proc interface{}) {
	sess.msgHdlr.RegMsgProc(proc)
}

// readAndHandleOneMsg readAndHandleOneMsg
func (sess *Session) readAndHandleOneMsg() error {
	var msgID inet.MsgID
	var rawMsgBuf []byte
	var err error

	if sess.isIdip {
		msgID, rawMsgBuf, err = readARQIdipMsg(sess.conn)
		//收到心跳消息
		if msgID == 0 && err == nil {
			sess.sendServerHB()
			return nil
		}
	} else {
		msgID, rawMsgBuf, err = readARQMsg(sess.conn, sess.isClient)
	}

	//log.Debug("readAndHandleOneMsg, msgID: ", msgID)

	if err != nil {
		return err
	}

	// 流量限制, 等待直到允许接收
	if err = sess.qpsLimiter.Wait(context.Background()); err != nil {
		return err
	}
	if err = sess.bpsLimiter.WaitN(context.Background(), len(rawMsgBuf)); err != nil {
		return err
	}

	assert.True(rawMsgBuf != nil, "rawMsg is nil")
	if !sess.needVerify {
		// 重置心跳計時器
		sess.ResetHb()
		sess.msgHdlr.HandleRawMsg(msgID, rawMsgBuf)
		return nil
	}

	// 会话需要验证，第1个消息为验证请求消息
	if msgID != sess.verifyMsgID {
		if sess.isIdip {
			return fmt.Errorf("need verify message ID %d, but got %d",
				sess.verifyMsgID, msgID)
		}
		return fmt.Errorf("need verify message ID %d, but got %d",
			sess.verifyMsgID, msgID)

	}
	// 重置心跳計時器
	sess.ResetHb()
	sess.msgHdlr.HandleRawMsg(msgID, rawMsgBuf)

	if sess.isVerified.Load() {
		sess.needVerify = false // 已通过验证，不再需要了，
		return nil
	}

	return fmt.Errorf("Session verification failed")
}

// SetVerified 设置会话已通过验证.
// thread-safe.
func (sess *Session) SetVerified() {
	sess.isVerified.Store(true)
}

// Emit Emit
func (sess *Session) Emit(evt string, p interface{}) {
	if f, ok := sess.evtMap[evt]; ok {
		f(p)
	}
}

// On On
func (sess *Session) On(evt string, f func(interface{})) {
	sess.evtMap[evt] = f
}

// SetVerifyMsgID 设置会话验证消息id
// 非线程安全，在Session.Start()之前设置。
func (sess *Session) SetVerifyMsgID(verifyMsgID inet.MsgID) {
	sess.verifyMsgID = verifyMsgID
	if sess.verifyMsgID != 0 {
		sess.needVerify = true
	}
}

// SetOnClosed SetOnClosed
func (sess *Session) SetOnClosed(onClosed func()) {
	sess.onClosed = onClosed
}

// SetBytePerSecLimiter 设置每秒接收字节数限制.
// r(rate) 为每秒字节数。
// b(burst) 为峰值字节数。
// 必须在 Start() 之前设置，避免 DataRace.
func (sess *Session) SetBytePerSecLimiter(r rate.Limit, b int) {
	sess.bpsLimiter = rate.NewLimiter(r, b)
}

// SetQueryPerSecLimiter 设置每秒接收请求数限制.
// r(rate) 为每秒请求数。
// b(burst) 为峰值请求数。
// 必须在 Start() 之前设置，避免 DataRace.
func (sess *Session) SetQueryPerSecLimiter(r rate.Limit, b int) {
	sess.qpsLimiter = rate.NewLimiter(r, b)
}

// AddOnClosed AddOnClosed
func (sess *Session) AddOnClosed(f func()) {
	sess.onClosedFuncs = append(sess.onClosedFuncs, f)
}

// GetID GetID
func (sess *Session) GetID() uint64 {
	return sess.id
}
