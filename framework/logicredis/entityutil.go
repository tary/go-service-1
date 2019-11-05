package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	dbservice "github.com/GA-TECH-SERVER/zeus/base/redisservice"
)

// EntityUtil 以hash表存储entity数据
// key: type:dbid
//如 player:1000
type EntityUtil struct {
	typ  string
	dbid uint64
}

// GetEntityUtil 获得Entity工具类
func GetEntityUtil(typ string, dbid uint64) *EntityUtil {
	enUtil := &EntityUtil{}
	enUtil.typ = typ
	enUtil.dbid = dbid
	return enUtil
}

// GetValues 获取值
func (util *EntityUtil) GetValues(args []interface{}) ([]interface{}, error) {
	c := dbservice.GetCacheConn()
	defer c.Close()

	return redis.Values(c.Do("HMGET", append([]interface{}{util.key()}, args...)...))
}

// SetValues 设置值
func (util *EntityUtil) SetValues(args []interface{}) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("HMSET", append([]interface{}{util.key()}, args...)...)
	return err
}

// GetValue 获取值
func (util *EntityUtil) GetValue(k string) (interface{}, error) {
	c := dbservice.GetCacheConn()
	defer c.Close()

	return c.Do("HGET", util.key(), k)
}

// SetValue 设置值
func (util *EntityUtil) SetValue(k string, v interface{}) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("HSET", util.key(), k, v)
	return err
}

// key key
func (util *EntityUtil) key() string {
	return fmt.Sprintf("%s:%d", util.typ, util.dbid)
}
