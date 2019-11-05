package userbase

import "github.com/giant-tech/go-service/base/net/inet"

// UserInitData 玩家初始化数据
type UserInitData struct {
	Sess    inet.ISession // session
	Version string        // 客户端版本
}
