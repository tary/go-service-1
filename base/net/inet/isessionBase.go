package inet

// ISessionBase session基础
type ISessionBase interface {
	Send(IMsg)
	//有其它需求再加
}
