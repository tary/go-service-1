package app

import (
	"github.com/giant-tech/go-service/framework/app/internal"
)

// Run 逻辑入口
// configFile 配置文件，如果为空字符串,就是从配置文件读取
func Run(configFile string) {
	internal.MyApp.Run(configFile)
}
