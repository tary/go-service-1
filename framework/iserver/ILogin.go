package iserver

import "github.com/giant-tech/go-service/framework/net/inet"

// ILogin 登录接口
type ILogin interface {
	OnLogin(uint64, inet.ISession) (IEntity, error)
}

// ILogout 登出接口
type ILogout interface {
	OnLogout(uint64)
}
