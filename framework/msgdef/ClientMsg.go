package msgdef

/*
	msgdef 包下面主要定义二进制消息流
	由于二进制消息流个数固定，所以在代码中直接定义，序列化和反序列化
*/

//ReturnType 返回类型
type ReturnType int32

const (
	//ReturnTypeSuccess 成功
	ReturnTypeSuccess ReturnType = 0
	//ReturnTypeServerBusy 服务器忙
	ReturnTypeServerBusy ReturnType = 1
	//ReturnTypeTokenInvalid token无效
	ReturnTypeTokenInvalid ReturnType = 2
	//ReturnTypeTokenParamErr 参数错误
	ReturnTypeTokenParamErr ReturnType = 3
	//ReturnTypeFailRelogin 重复登录
	ReturnTypeFailRelogin ReturnType = 4
	//ReturnTypeFailWrongVersion 版本号错误
	ReturnTypeFailWrongVersion ReturnType = 5
)

// LoginReq 登录请求
// Client ==> LobbyServer
type LoginReq struct {
	Account string // 账号
	Token   string // token
	UID     uint64 // 玩家ID
	Version string // 版本
	ExtData []byte // 自定义数据
}

// LoginResp 登录返回
// LobbyServer ==> Client
type LoginResp struct {
	Result  uint32 // 返回类型
	ErrStr  string // 错误内容
	ExtData []byte // 自定义数据
}

// Ping ping
type Ping struct {
}

// Pong pong
type Pong struct {
}

// CreateEntityNotify 创建实体通知
type CreateEntityNotify struct {
	EntityType string
	EntityID   uint64
}

// CallMsg 远程调用消息
type CallMsg struct {
	GroupID      uint64 // 目标实体所属的GroupID
	EntityID     uint64 // 目的EntityID
	SType        uint8  // 目标服务类型
	SID          uint64 // 目标服务ID
	FromSID      uint64 // From Service ID
	Seq          uint64 // 序号
	MethodName   string // 方法名
	Params       []byte // 参数
	IsSync       bool   // 是否为同步
	IsFromClient bool   // 是否来自客户端
}

// CallRespMsg 远程调用的返回消息
type CallRespMsg struct {
	Seq       uint64 // 序号
	ErrString string // 错误
	RetData   []byte // 返回的数据
}

// ForwardToClientMsg 转发给客户端
type ForwardToClientMsg struct {
	ServiceID uint64
	EntityID  uint64
	MsgData   []byte
}
