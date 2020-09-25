package inet

// ISDSession session基础 主要用于service discovery的session
type ISDSession interface {
	Send(IMsg) error
	//有其它需求再加
}
