package userbase

import "github.com/GA-TECH-SERVER/zeus/base/net/inet"

// UserInitData 玩家初始化数据
type UserInitData struct {
	Sess    inet.ISession // session
	Version string        // 客户端版本
}
