package iserver

//ICall call接口
type ICall interface {
	PostAct(f func())
	Call(f func() interface{}) interface{}
}
