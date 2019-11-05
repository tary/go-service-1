package login

// UserLoginReq 登录消息格式
type UserLoginReq struct {
	User     string
	Password string
	Channel  string
	Data     []byte
	GameType uint32
}

// UserLoginAck 登录消息返回格式
type UserLoginAck struct {
	UID       uint64
	Token     string
	LobbyAddr string
	Result    int
	ResultMsg string
	HB        bool
	// Config    interface{}
}

// LVerifyReq 登录Token验证
type LVerifyReq struct {
	UID   uint64
	Token string
}

// LVerifyAck 登录Token验证结果
type LVerifyAck struct {
	Result int
}

// UpdateLoadReq 向Login报告负载
type UpdateLoadReq struct {
	OuterAddr string
	Load      int
	GameType  uint32
}

// UserLogoutReq 登出消息格式
type UserLogoutReq struct {
	UID       uint64
	OuterAddr string
}

/*
	Client <==> LobbyServer  HTTP
*/
// UserLoginReq 登录消息格式
/*type MatchReq struct {
	MatchType uint32 // 匹配的房间类型
	Number    uint32 // 当前房间内的匹配人数
	TickRate  uint32 // 创建房间的帧
	CloseRoom uint32 // 多少秒后关闭房间
}

// LobbyServer ==> Client
type MatchResp struct {
	Token    string // token值
	RoomIp   string // RelayServer地址
	RoomPort uint16 // RelayServer端口
}
*/

// CheckTokenReq LobbySrv < == > LoginSvr (http protocol)
type CheckTokenReq struct {
	UserID uint64 // 客户端ID
	Token  string // 客户端从LoginSrv中获取的Token
}

// CheckTokenResp < == > LobbySrv (http protocol)
type CheckTokenResp struct {
	Ok     bool   // LoginSrv Token检查的返回结果
	Result string // 不成功情况下返回的错误信息
}
