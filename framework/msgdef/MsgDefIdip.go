package msgdef

import (
	"reflect"

	"github.com/giant-tech/go-service/base/net/inet"
)

// 基础消息id
const (
	MsgIDServerHB   = 0   // 服务器心跳
	MsgIDClientHB   = 1   // 客户端心跳
	MsMsgIDLogin    = 257 // 登录0x101
	MsMsgIDLoginOK  = 513 // 登录成功0x201
	MsgIDClientInfo = 870 // 客户端信息
	MsgIDChatCmd    = 650 // 聊天消息
	MsgIDIdipCmd    = 614 // IdIp消息  40056
)

const (
	// SERVICENAMELENGTH 登录0x101
	SERVICENAMELENGTH = 16
	// AUTHENTICATELENGTH 登录成功0x201
	AUTHENTICATELENGTH = 32
	// ERRORMSGLENGTH 错误码
	ERRORMSGLENGTH = 100
	// IDIPBODYLENGTH body len
	IDIPBODYLENGTH = 24000
)

// GameZone 游戏区
type GameZone struct {
	ZoneID uint16
	GameID uint16
}

// LoginCmd 登录命令
type LoginCmd struct {
	IP     [16]byte
	Port   uint16
	GameID uint16
}

// LoginCmdOK 登录命令ok
type LoginCmdOK struct {
	GameZone GameZone
	Name     [32]byte
	NetType  byte
}

// ChatCmd 聊天命令
type ChatCmd struct {
	CmdID uint32
	Data  []byte
}

// ClientInfo client info
type ClientInfo struct {
	ZoneID        uint16
	ClientVersion [65]byte
}

// IdipHeader idip头
type IdipHeader struct {
	PacketLen    uint32                   /* 包体长度 */
	Cmdid        uint32                   /* 命令ID */
	Seqid        uint32                   /* 流水号 */
	ServiceName  [SERVICENAMELENGTH]byte  /* 服务名 */
	SendTime     uint32                   /* 发送时间YYYYMMDD对应的整数 */
	Version      uint32                   /* 版本号 */
	Authenticate [AUTHENTICATELENGTH]byte /* 加密串 */
	Result       int32
	/* 错误码,返回码类型：0：处理成功，需要解开包体获得详细信息,1：处理成功，但包体返回为空，不需要处理包体（eg：查询用户角色，用户角色不存在等），-1: 网络通信异常,-2：超时,-3：数据库操作异常,-4：API返回异常,-5：服务器忙,-6：其他错误,小于-100 ：用户自定义错误，需要填写szRetErrMsg */
	RetErrMsg [ERRORMSGLENGTH]byte /* 错误信息 */
}

// IdipCmd idip cmd
type IdipCmd struct {
	HTTPID      uint32
	EventID     uint64
	IsBroadcast byte
	IdipHead    IdipHeader
	Data        []byte
}

// ClientHB 客户端主动发起的心跳
type ClientHB struct {
}

// ServerHB 服务器主动发起的心跳
type ServerHB struct {
}

// InitIDIP 初始消息定义管理器
func (def *MsgDef) InitIDIP() {
	def.RegMsg(MsgIDServerHB, &ServerHB{})
	def.RegMsg(MsgIDClientHB, &ClientHB{})
	def.RegMsg(MsMsgIDLogin, &LoginCmd{})
	def.RegMsg(MsMsgIDLoginOK, &LoginCmdOK{})
	def.RegMsg(MsgIDClientInfo, &ClientInfo{})
	def.RegMsg(MsgIDChatCmd, &ChatCmd{})
	def.RegMsg(MsgIDIdipCmd, &IdipCmd{})
}

// msgDefInstIdip 实例
var msgDefInstIdip *MsgDef

// GetMsgDefIDIP 获取消息定义对象的全局实例
func GetMsgDefIDIP() *MsgDef {
	if msgDefInstIdip == nil {
		msgDefInstIdip = &MsgDef{
			make(map[inet.MsgID]*MsgInfo),
			make(map[string]*MsgInfo),
			make(map[reflect.Type]*MsgInfo),
		}
		msgDefInstIdip.InitIDIP()
	}

	return msgDefInstIdip
}

// InitIDIP 初始化
func InitIDIP() {
	if msgDefInstIdip == nil {
		msgDefInstIdip = &MsgDef{
			make(map[inet.MsgID]*MsgInfo),
			make(map[string]*MsgInfo),
			make(map[reflect.Type]*MsgInfo),
		}
		msgDefInstIdip.InitIDIP()
	}
}
