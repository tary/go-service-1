package igateway

import (
	"github.com/GA-TECH-SERVER/zeus/base/net/inet"
	"github.com/GA-TECH-SERVER/zeus/framework/msgdef"
)

// ReconnectData 断线重连返回的数据
type ReconnectData struct {
	IsCreateEntity bool  // 是否创建新实体，true说明当前玩家会被踢掉外面再创建一个新对象，false说明继续使用当前对象
	Err            error // 错误内容
}

// ILoginHandler 登录处理
type ILoginHandler interface {
	OnLoginHandler(msg *msgdef.LoginReq) msgdef.ReturnType
}

// IReconnectHandler 重连处理
type IReconnectHandler interface {
	// OnReconnect 重连，可能是断线重连可能是顶号
	OnReconnect(sess inet.ISession) *ReconnectData
}

// ICloseHandler 登录处理
type ICloseHandler interface {
	OnClose()
}

// IGateUser 接口
type IGateUser interface {
	OnUserInit() error
	OnUserTick()
	OnUserFini()

	IReconnectHandler
	ICloseHandler
}
