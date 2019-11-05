package dbservice

import (
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// EntityDBIDUtil 以hash表存储entity数据
// key: entitydbid:dbid
type EntityDBIDUtil struct {
	dbid uint64
}

const (
	// entityDBIDPrefix entityDBIDPrefix
	entityDBIDPrefix = "entitydbid"
	// fieldEntityID fieldEntityID
	fieldEntityID = "entityid"
	// fieldExistedCount fieldExistedCount
	fieldExistedCount = "existedcnt"
)

// GetEntityDBIDUtil Entity DBID工具类
func GetEntityDBIDUtil(dbid uint64) *EntityDBIDUtil {
	enUtil := &EntityDBIDUtil{}
	enUtil.dbid = dbid
	return enUtil
}

// IsExist 这个ID号的Entity是否存在
func (util *EntityDBIDUtil) IsExist() bool {
	c := dbservice.GetCacheConn()
	defer c.Close()

	r, err := c.Do("EXISTS", util.key())
	if err != nil {
		log.Error(err)
		return false
	}

	v, err := redis.Bool(r, nil)
	if err != nil {
		log.Error(err)
		return false
	}
	return v
}

// RegEntityID 注册类型和服务器信息
func (util *EntityDBIDUtil) RegEntityID(entityID uint64) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("HMSET", util.key(), fieldEntityID, entityID)
	if err != nil {
		log.Error(err)
	}
	_, err = c.Do("HINCRBY", util.key(), fieldExistedCount, 1)
	return err
}

// UnRegEntityID 删除注册信息
func (util *EntityDBIDUtil) UnRegEntityID() error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	reply, err := c.Do("HINCRBY", util.key(), fieldExistedCount, -1)
	if err != nil {
		log.Error(err)
	}

	cnt, err := redis.Int(reply, nil)
	if err != nil {
		log.Error(err)
	}
	if cnt <= 0 {
		_, err = c.Do("DEL", util.key())
		return err
	}

	_, err = c.Do("HDEL", util.key(), fieldEntityID)
	if err != nil {
		log.Error(err)
	}

	return nil
}

// GetEntityID 获取entityID
func (util *EntityDBIDUtil) GetEntityID() (uint64, error) {
	c := dbservice.GetCacheConn()
	defer c.Close()

	ret, err := c.Do("HMGET", util.key(), fieldEntityID)
	if err != nil {
		return 0, err
	}

	retValue, err := redis.Uint64(redis.Values(ret, nil))
	if err != nil {
		return 0, err
	}

	return retValue, nil
}

// key key
func (util *EntityDBIDUtil) key() string {
	return fmt.Sprintf("%s:%d", entityDBIDPrefix, util.dbid)
}
