package gatewaybase

// MatchReq 匹配请求
type MatchReq struct {
	MatchKey string //匹配key
}

// CancelMatchReq 取消匹配请求
type CancelMatchReq struct {
	MatchKey string //匹配key
}

// MatchResp 匹配回应
type MatchResp struct {
	ReturnType int32
	ExpectTime uint64
}
