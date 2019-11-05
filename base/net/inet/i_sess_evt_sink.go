package inet

// ISessEvtSink sess处理
type ISessEvtSink interface {
	OnConnected(ISession)
	OnClosed(ISession)
}
