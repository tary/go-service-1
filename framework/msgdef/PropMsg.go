package msgdef

// SyncProps 同步属性
type SyncProps struct {
	EntityID uint64
	PropNum  uint32
	Data     []byte
}

// PlayerProp 玩家属性
type PlayerProp struct {
	PropName  string `protobuf:"bytes,1,req,name=propName" json:"propName"`
	PropValue []byte `protobuf:"bytes,2,req,name=propValue" json:"propValue"`
}

// PlayerProps 玩家属性结构体数组
type PlayerProps struct {
	Props []*PlayerProp `protobuf:"bytes,1,rep,name=props" json:"props,omitempty"`
}
