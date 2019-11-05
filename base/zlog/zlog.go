package zlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitDefault 初始化日志库
func InitDefault() {
	logFileName := viper.GetString("Log.LogFileName")
	logConfig := viper.GetString("Log.LogConfig")
	dir := viper.GetString("Log.LogDir")
	Init(dir, logFileName, logConfig)
}

// Init 初始化日志库
// logname 为日志文件名， 如果为空则使用可执行文件名
// logConfig 为日志配置文件
func Init(dir, logFileName, logConfig string) {
	load(dir, logFileName, logConfig)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("zlog panic:", err, ", Stack: ", string(debug.Stack()))
				if viper.GetString("Config.Recover") == "0" {
					panic(err)
				}
			}
		}()

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error(err)
			return
		}
		defer watcher.Close()

		if err = watcher.Add(logConfig); err != nil {
			log.Error(err)
			return
		}

		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op == fsnotify.Write {
					// if err := viper.ReadInConfig(); err != nil {
					// 	log.Error(err)
					// 	continue
					// }

					load(dir, logFileName, logConfig)
				}
			case err := <-watcher.Errors:
				log.Error(err)
			}
		}
	}()
}

func load(logdir, logFileName, logConfig string) {
	fmt.Println("Logdir:", logdir, ", logFileName:", logFileName, ", logConfig: ", logConfig)

	b, err := ioutil.ReadFile(logConfig)
	if err != nil {
		//fmt.Print(err)
		log.Error("read logConfig file err:", err)
		return
	}

	var srvName string

	if len(logFileName) == 0 {
		srvName = logdir + filepath.Base(os.Args[0])
	} else {
		srvName = logdir + logFileName
	}

	oldFile := string(b)
	configFile := strings.Replace(oldFile, "tempLogFileName", srvName, -1)

	defer log.Flush()
	logger, err := log.LoggerFromConfigAsString(configFile)
	if err != nil {
		panic(err)
	}

	//logger.SetAdditionalStackDepth(1)

	log.ReplaceLogger(logger)
}
