package msgdef

const (
	//LoginReqMsgID 登录请求
	LoginReqMsgID = 1
	//LoginRespMsgID 登录回应
	LoginRespMsgID = 2
	//PingMsgID ping消息
	PingMsgID = 3
	//PongMsgID pong消息
	PongMsgID = 4
	//CallMsgID 函数调用
	CallMsgID = 5
	//CallRespMsgID 函数调用返回
	CallRespMsgID = 6
	//CreateEntityNotifyMsgID 创建实体通知
	CreateEntityNotifyMsgID = 7

	//ForwardToClientMsgID 转发给客户端消息
	ForwardToClientMsgID = 8

	//EntitySrvInfoNotifyID 实体变化消息
	EntitySrvInfoNotifyID = 15
	//ProtoSyncMsgID ProtoSync消息ID
	ProtoSyncMsgID = 16
	// ClientFrameMsgID 客户端过来的帧消息
	ClientFrameMsgID = 20
	//ServerFrameMsgID 服务器下发的单帧消息
	ServerFrameMsgID = 21
	//FramesMsgID 服务器下发的多帧消息
	FramesMsgID = 22
	//RequireFramesMsgID 客户端
	RequireFramesMsgID = 23
	//UserDuplicateLoginNotifyMsgID 玩家重复登录的消息
	UserDuplicateLoginNotifyMsgID = 24

	//PropsSyncMsgID 属性消息ID
	PropsSyncMsgID = 30
	// PropsSyncClientMsgID 发给客户端的属性消息ID
	PropsSyncClientMsgID = 31
	//MRolePropsSyncClientMsgID 主角属性同步消息DI
	MRolePropsSyncClientMsgID = 32

	//EnterCellReqMsgID 进入空间的消息
	EnterCellReqMsgID = 40
	//LeaveCellReqMsgID 离开空间的消息
	LeaveCellReqMsgID = 41
	//EnterAOIMsgID 进入AOI消息
	EnterAOIMsgID = 42
	//LeaveAOIMsgID 离开AOI消息
	LeaveAOIMsgID = 43
	//AOIPosChangeMsgID AOI位置改变消息
	AOIPosChangeMsgID = 44
	//EnterCellMsgID 玩家进入场景
	EnterCellMsgID = 45
	//LeaveCellMsgID 玩家离开场景
	LeaveCellMsgID = 46
	//UserMoveMsgID 玩家移动
	UserMoveMsgID = 47
	//CellEntityMsgID 空间实体消息
	CellEntityMsgID = 48
	//EntityPosSetMsgID 玩家移动
	EntityPosSetMsgID = 49

	//ClientTransportMsgID 客户端中转消息ID
	ClientTransportMsgID = 50
	//CreateBaseReqMsgID 请求创建base实体消息ID
	CreateBaseReqMsgID = 51
	//CreateBaseRetMsgID 创建base实体返回消息ID
	CreateBaseRetMsgID = 52
	//DestroyBaseReqMsgID 请求删除base实体消息ID
	DestroyBaseReqMsgID = 53
	//DestroyBaseRetMsgID 销毁base实体返回消息ID
	DestroyBaseRetMsgID = 54
	// CreateCellEntityReqMsgID  CreateCellEntityReqMsgID
	CreateCellEntityReqMsgID = 55
	//CreateCellEntityRetMsgID 创建base实体返回消息ID
	CreateCellEntityRetMsgID = 56
	//DestroyCellEntityReqMsgID 请求删除base实体消息ID
	DestroyCellEntityReqMsgID = 57
	// DestroyCellEntityRetMsgID 销毁base实体返回消息ID
	DestroyCellEntityRetMsgID = 58

	//EntityMsgTransportMsgID 分布式实体之间传递消息用
	EntityMsgTransportMsgID = 59
	//EntityMsgChangeMsgID 分布式实体之间同步数据使用
	EntityMsgChangeMsgID = 60
	//SrvMsgTransportMsgID 服务器间消息转发ID
	SrvMsgTransportMsgID = 61
	//RPCMsgID RPC消息
	RPCMsgID = 62
	//SyncClockMsgID SyncClockMsgID
	SyncClockMsgID = 63
	//SpaceUserConnectMsgID SpaceUserConnectMsgID
	SpaceUserConnectMsgID = 71
	//SpaceUserConnectSucceedRetMsgID SpaceUserConnectSucceedRetMsgID
	SpaceUserConnectSucceedRetMsgID = 72
	// SyncUserStateMsgID SyncUserStateMsgID
	SyncUserStateMsgID = 73
	// AOISyncUserStateMsgID AOISyncUserStateMsgID
	AOISyncUserStateMsgID = 74
	//AdjustUserStateMsgID AdjustUserStateMsgID
	AdjustUserStateMsgID = 75
	//EntityAOISMsgID EntityAOISMsgID
	EntityAOISMsgID = 76
	//EntityBasePropsMsgID EntityBasePropsMsgID
	EntityBasePropsMsgID = 77
	//EntityEventMsgID EntityEventMsgID
	EntityEventMsgID = 78

	//CellMsgTransportMsgID 发消息给某个cell
	CellMsgTransportMsgID = 81

	//CreateGhostReqMsgID 创建ghost
	CreateGhostReqMsgID = 82
	//DeleteGhostReqMsgID 删除ghost
	DeleteGhostReqMsgID = 83
	//TransferRealReqMsgID 切换real
	TransferRealReqMsgID = 84
	//NewRealNotifyMsgID 新real通知
	NewRealNotifyMsgID = 85
	//VideoListReqMsgID VideoListReqMsgID
	VideoListReqMsgID = 106
	//WatchVideoReqMsgID WatchVideoReqMsgID
	WatchVideoReqMsgID = 107
	//CancelWatchReqMsgID CancelWatchReqMsgID
	CancelWatchReqMsgID = 108
	//BroadcastRoomMsgMsgID BroadcastRoomMsgMsgID
	BroadcastRoomMsgMsgID = 109
	//PlayerGmDebugRespMsgID PlayerGmDebugRespMsgID
	PlayerGmDebugRespMsgID = 110
	//GameEndReqMsgID GameEndReqMsgID
	GameEndReqMsgID = 111
	//MatchReqMsgID MatchReqMsgID
	MatchReqMsgID = 112
	//ForwardMsgToClientMsgID ForwardMsgToClientMsgID
	ForwardMsgToClientMsgID = 113
	//ErrorMessageMsgID ErrorMessageMsgID
	ErrorMessageMsgID = 114
	//PlayerGmDebugReqMsgID PlayerGmDebugReqMsgID
	PlayerGmDebugReqMsgID = 115

	//CancelMatchReqMsgID CancelMatchReqMsgID
	CancelMatchReqMsgID = 117
	//PlayerDisconnectNotifyMsgID PlayerDisconnectNotifyMsgID
	PlayerDisconnectNotifyMsgID = 118
	//PlayerGmRespMsgID PlayerGmRespMsgID
	PlayerGmRespMsgID = 120
	//SyncPropsMsgID SyncPropsMsgID
	SyncPropsMsgID = 121
	//GameEndReqCenterMsgID GameEndReqCenterMsgID
	GameEndReqCenterMsgID = 122

	//PlayerGmReqMsgID PlayerGmReqMsgID
	PlayerGmReqMsgID = 124
	//MatchRespMsgID MatchRespMsgID
	MatchRespMsgID = 125
	//MatchSuccessMsgID MatchSuccessMsgID
	MatchSuccessMsgID = 126
	//MatchClientReqMsgID MatchClientReqMsgID
	MatchClientReqMsgID = 127
	//MatchClientRespMsgID MatchClientRespMsgID
	MatchClientRespMsgID = 128
	//ClientVerifyReqMsgID ClientVerifyReqMsgID
	ClientVerifyReqMsgID = 129
	//ClientVerifyRespMsgID ClientVerifyRespMsgID
	ClientVerifyRespMsgID = 130

	// EnterReqMsgID EnterReqMsgID  relay消息
	EnterReqMsgID = 141
	// EnterRespMsgID EnterRespMsgID
	EnterRespMsgID = 142
	//ActMsgMsgID ActMsgMsgID
	ActMsgMsgID = 143
	//TickDataMsgID TickDataMsgID
	TickDataMsgID = 144

	//StartMatchReqMsgID StartMatchReqMsgID MatchServer消息
	StartMatchReqMsgID = 150
	//StartMatchRespMsgID StartMatchRespMsgID
	StartMatchRespMsgID = 151
	//StopMatchReqMsgID StopMatchReqMsgID
	StopMatchReqMsgID = 152
	//StopMatchRespMsgID StopMatchRespMsgID
	StopMatchRespMsgID = 153
	//CreateMatchRoomNotifyMsgID CreateMatchRoomNotifyMsgID
	CreateMatchRoomNotifyMsgID = 154
	//LeaveMatchRoomNotifyMsgID LeaveMatchRoomNotifyMsgID
	LeaveMatchRoomNotifyMsgID = 155
	//JoinMatchRoomNotifyMsgID JoinMatchRoomNotifyMsgID
	JoinMatchRoomNotifyMsgID = 156
	//MatchFinishNotifyMsgID MatchFinishNotifyMsgID
	MatchFinishNotifyMsgID = 157
)

const (
	// ClientMSG 来自客户端的验证消息
	ClientMSG uint8 = 0
)

// ClientVertifySucceedRet 中登录类型
const (
	// Connected 正常连接
	Connected uint8 = 1
	// ReConnect 断线重连
	ReConnect uint8 = 2
	// DupConnect 重复连接
	DupConnect uint8 = 3
)

const (
	// SessHBTimeout 心跳超时
	SessHBTimeout uint8 = 1
	// SessDisconnect 断线
	SessDisconnect uint8 = 2
	// SessError 出错
	SessError uint8 = 3
)
