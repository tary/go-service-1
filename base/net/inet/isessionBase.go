package inet

// ISessionBase session基础
type ISessionBase interface {
	Send(IMsg) error
	//有其它需求再加
}
