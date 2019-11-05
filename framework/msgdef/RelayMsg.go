package msgdef

// RoomEventType 房间事件枚举
type RoomEventType int32

const (
	// RoomEventTypeRETNONE none
	RoomEventTypeRETNONE RoomEventType = 0
	// RoomEventTypePlayerOnline online
	RoomEventTypePlayerOnline RoomEventType = 1
	// RoomEventTypePlayerOffline offline
	RoomEventTypePlayerOffline RoomEventType = 2
	// RoomEventTypeRoomClosed closed
	RoomEventTypeRoomClosed RoomEventType = 3
)

// RoomEvent 房间事件
type RoomEvent struct {
	EventType RoomEventType
	PlayerID  uint32
}

// EnterReq 进入请求
type EnterReq struct {
	// 服务器可以从 player_token 找到对应的房间和玩家号.
	PlayerToken string
	// 上个帧号，需要服务器发送之后的帧: last_tick+1, last_tick+2, ...
	// 用于断线重连。
	// 如果是从头开始，则 last_tick 为0, 开始帧为 1.
	LastTick uint32
}

// EnterResp 进入回应
type EnterResp struct {
	Ok       bool
	Error    string
	PlayerID uint32
	Ifs      uint32
	RandSeed uint64
	RoomName string
}

// TickMsgPlayerActions 单个玩家当前帧的所有输入动作。保持输入次序。
type TickMsgPlayerActions struct {
	PlayerID uint32
	Actions  [][]byte
}

// TickData 帧数据。
// 中继服将记录所有帧消息，用于断线重连。
type TickData struct {
	Tick          uint32
	PlayerActions []*TickMsgPlayerActions
	Events        []*RoomEvent
}

// ActMsg 动作消息
type ActMsg struct {
	// 客户端自定义格式，服务器不会解析.
	// 多个动作一起发送。
	Actions [][]byte
}

// CardData 可选卡牌数据
type CardData struct {
	CardID uint32 `protobuf:"varint,1,req,name=cardID" json:"cardID"`
}

// PlayerExtData 玩家比赛用到的额外信息
type PlayerExtData struct {
	Cards []*CardData `protobuf:"bytes,1,rep,name=cards" json:"cards,omitempty"`
}

// VideoListReq 录像请求
type VideoListReq struct {
	Tp uint32
}

// WatchVideoReq 查看录像请求
type WatchVideoReq struct {
	RoomName string
	PlayerID uint32
}

// CancelWatchReq 取消查看请求
type CancelWatchReq struct {
}
