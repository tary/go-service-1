package app

import (
	"github.com/GA-TECH-SERVER/zeus/framework/app/internal"
)

// Run 逻辑入口
// configFile 配置文件，如果为空字符串
func Run(configFile string) {
	internal.MyApp.Run(configFile)
}
