package iserver

import "github.com/GA-TECH-SERVER/zeus/base/net/inet"

// ILogin 登录接口
type ILogin interface {
	OnLogin(uint64, inet.ISession) (IEntity, error)
}

// ILogout 登出接口
type ILogout interface {
	OnLogout(uint64)
}
