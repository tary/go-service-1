package sbase

//server base package
//解决包循环依赖

//loginMgr 服务器级别最低的包，供其他包使用
var loginMgr ILoginMgr

// ILoginMgr ILoginMgr接口
type ILoginMgr interface {
	PlayerLogout(uint64) error
}

// SetLoginMgr 设置登录管理器
func SetLoginMgr(icm ILoginMgr) {
	loginMgr = icm
}

// GetLoginMgr 得到登录管理器
func GetLoginMgr() ILoginMgr {
	return loginMgr
}
