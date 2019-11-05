package consts

const (
	// MsgIDSize id size
	MsgIDSize = 2
	// MsgHeadSize consist message length , compression type
	// 3字节长度+1字节标记位
	MsgHeadSize = 4

	// IdipMsgIDSize idip size
	IdipMsgIDSize = 2
	// IdipMsgHeadSize idip消息 2字节标志，2字节长度，2字节msgID
	IdipMsgHeadSize = 6
)

const (
	// MaxMsgBuffer 消息最大长度
	MaxMsgBuffer = 100 * 1024
)
