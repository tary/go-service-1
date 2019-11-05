package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// CellUtil cell相关
type CellUtil struct {
	id uint64
}

const (
	cellPrefix = "cell"
)

// GetCellUtil 获取Util
func GetCellUtil(id uint64) *CellUtil {
	return &CellUtil{id: id}
}

// RegSrvID 设置服务器ID
func (util *CellUtil) RegSrvID(srvID uint64) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("HSET", util.key(), "SrvID", srvID)
	if err != nil {
		return err
	}

	return nil
}

// UnReg 设置服务器
func (util *CellUtil) UnReg() error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("DEL", util.key())
	if err != nil {
		return err
	}
	return nil
}

// IsExist 是否存在
func (util *CellUtil) IsExist() (bool, error) {

	c := dbservice.GetCacheConn()
	defer c.Close()

	r, err := c.Do("EXISTS", util.key())
	if err != nil {
		return false, err
	}

	ret, err := redis.Bool(r, nil)
	if err != nil {
		return false, err
	}

	return ret, nil
}

// GetSrvID 获取服务器ID
func (util *CellUtil) GetSrvID() (uint64, error) {

	c := dbservice.GetCacheConn()
	defer c.Close()

	r, err := c.Do("HGET", util.key(), "SrvID")
	if err != nil {
		return 0, err
	}

	srvID, err := redis.Uint64(r, nil)
	if err != nil {
		return 0, err
	}

	return srvID, nil
}

// key key
func (util *CellUtil) key() string {
	return fmt.Sprintf("%s:%d", cellPrefix, util.id)
}
