package dbservice

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	dbservice "github.com/GA-TECH-SERVER/zeus/base/redisservice"
	"github.com/GA-TECH-SERVER/zeus/framework/idata"
)

/*
 server:* 表工具类
*/

const (
	servicePrefix = "AppService:"
	serviceIDs    = "AppServiceIDs"
	serviceInfo   = "AppServiceInfo"
)

// ServiceUtil service相关
type ServiceUtil struct {
	svrID uint64
}

// GetServiceUtil 获取service相关
func GetServiceUtil(svrID uint64) *ServiceUtil {
	srvutil := &ServiceUtil{}
	srvutil.svrID = svrID
	return srvutil
}

// GetServiceList 获得服务列表
func GetServiceList(slist *[]*idata.ServiceInfo) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	ids, err := redis.Ints(c.Do("SMEMBERS", serviceIDs))
	if err != nil {
		return err
	}

	for _, id := range ids {
		key := fmt.Sprintf("%s%d", servicePrefix, id)
		v, err := c.Do("EXISTS", key)
		if err != nil {
			return err
		}
		if v.(int64) == 0 {
			_, err := c.Do("SREM", serviceIDs, id)
			if err != nil {
				return err
			}
			continue
		}

		data, err := redis.Bytes(c.Do("HGET", key, serviceInfo))
		if err != nil {
			continue
		}

		sInfo := &idata.ServiceInfo{}
		err = json.Unmarshal(data, sInfo)
		if err != nil {
			return err
		}

		*slist = append(*slist, sInfo)
	}

	return nil
}

// SetStatus 设置服务状态
func (util *ServiceUtil) SetStatus(status int) error {
	return util.setValue("status", status)
}

// Register 注册服务器信息, 设置过期时间30秒
func (util *ServiceUtil) Register(sInfo *idata.ServiceInfo) error {
	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()

	data, err := json.Marshal(sInfo)
	if err != nil {
		return err
	}

	if _, err := c.Do("HSET", util.key(), serviceInfo, data); err != nil {
		return err
	}
	if _, err := c.Do("SADD", serviceIDs, util.svrID); err != nil {
		return err
	}
	if _, err := c.Do("EXPIRE", util.key(), 30); err != nil {
		return err
	}
	return nil
}

// Update 更新服务器信息, 设置过期时间30秒
func (util *ServiceUtil) Update(sInfo *idata.ServiceInfo) error {
	data, err := json.Marshal(sInfo)
	if err != nil {
		return err
	}

	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()

	if _, err := c.Do("HSET", util.key(), serviceInfo, data); err != nil {
		return err
	}
	if _, err := c.Do("EXPIRE", util.key(), 30); err != nil {
		return err
	}
	return nil
}

// IsExist 当前服务器是否已经注册过
func (util *ServiceUtil) IsExist() bool {
	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()

	if isExist, err := redis.Bool(c.Do("EXISTS", util.key())); err == nil {
		return isExist
	}

	return false
}

// RefreshExpire 刷新过期时间
func (util *ServiceUtil) RefreshExpire(expire uint32) error {
	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("EXPIRE", util.key(), expire)
	return err
}

// Delete 删除key
func (util *ServiceUtil) Delete() error {
	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()
	if _, err := c.Do("DEL", util.key()); err != nil {
		return err
	}
	if _, err := c.Do("SREM", servers, util.svrID); err != nil {
		return err
	}
	return nil
}

func (util *ServiceUtil) setValue(field string, value interface{}) error {
	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()
	_, err := c.Do("HSET", util.key(), field, value)
	return err
}

func (util *ServiceUtil) getValue(field string) (interface{}, error) {
	//c := GetSingletonRedis()
	c := dbservice.GetCacheConn()
	defer c.Close()
	return c.Do("HGET", util.key(), field)
}

func (util *ServiceUtil) key() string {
	return fmt.Sprintf("%s%d", servicePrefix, util.svrID)
}
