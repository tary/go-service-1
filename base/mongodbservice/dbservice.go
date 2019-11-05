package dbservice

import (
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
	"go.uber.org/atomic" //"gopkg.in/mgo.v2"
)

var (
	// 仅初始化一次
	once sync.Once

	// isDBValid DB是否正常
	isDBValid = atomic.NewBool(true)

	globalMgoSession *mgo.Session /* = nil*/
)

//测试时这个函数再打开
func setConfig(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败")
	}
}

//InitDB 初始化db
func InitDB() bool {

	var err error
	var session *mgo.Session
	if globalMgoSession == nil {

		//测试时这两行再打开
		//configPath := "../../../res/config/server.toml"
		//setConfig(configPath)
		addr := viper.GetString("MongoDB.Addr")
		timeout := viper.GetInt64("MongoDB.Timeout")
		if timeout == 0 {
			timeout = 1
		}

		dur := time.Duration(timeout) * time.Second

		session, err = mgo.DialWithTimeout(addr, dur)

		if err != nil {
			log.Debug("connect failed,", err.Error())
			return false
			//panic(err)
		}
		log.Info("connect MongoDB success, addr = ", addr)
	}

	globalMgoSession = session
	globalMgoSession.SetMode(mgo.Monotonic, true)
	//default is 4096
	globalMgoSession.SetPoolLimit(300)
	return true
}

//CloneSession 克隆一个session
func CloneSession() *mgo.Session {
	if globalMgoSession == nil {
		rv := InitDB()
		if !rv {
			return nil
		}
	}
	return globalMgoSession.Clone()
}

/*func GetMongodbConn() *mgo.Session {
	once.Do(initMongoDB)
	return GlobalMgoSession
}*/
