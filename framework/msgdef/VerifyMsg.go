package msgdef

// ClientVerifyReq ClientVerifyReq
type ClientVerifyReq struct {
	ServerType uint32
	ServerID   uint64
	Token      string
}

// ClientVerifyResp ClientVerifyResp
type ClientVerifyResp struct {
	Result     uint32
	ServerType uint32
	ServerID   uint64
}
