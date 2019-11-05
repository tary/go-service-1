package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	dbservice "github.com/GA-TECH-SERVER/zeus/base/redisservice"
)

const (
	serverIDPrefix = "appserverid"
)

// SrvIDUtil srv相关
type SrvIDUtil struct {
	srvID uint64
}

// GetSrvIDUtil 获得SrvID工具类
func GetSrvIDUtil() *SrvIDUtil {
	return &SrvIDUtil{}
}

// GetServerID 获取ServerID并累加
func (util *SrvIDUtil) GetServerID() (uint64, error) {
	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()

	n, err := redis.Uint64(c.Do("INCRBY", util.key(), 1))
	if err != nil {
		return 0, err
	}

	return n, nil
}

// key key
func (util *SrvIDUtil) key() string {
	return fmt.Sprintf("%s", serverIDPrefix)
}
