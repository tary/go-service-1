package dbservice

import (
	"errors"
	"fmt"
	"time"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// defaultSrvName 暂仅支持一种服务类型，后期如有需要再扩展
const defaultSrvName = "GameServer"

// TIMEOUT TIMEOUT
const TIMEOUT = 30000 // (ms)

// SrvLoadUtil 相关
type SrvLoadUtil struct {
	srvName string // 服务名，如：demo_server
	srvAddr string // 服务地址
}

// GetSrvLoadUtil  SrvLoadUtil
func GetSrvLoadUtil(srvAddr string) *SrvLoadUtil {
	return &SrvLoadUtil{
		srvName: defaultSrvName,
		srvAddr: srvAddr,
	}
}

// GetMinLoadSrv 得到最小的服务器压力
func GetMinLoadSrv() (string, error) {

	c := dbservice.GetCacheConn()
	defer c.Close()

	values, err := redis.Strings(c.Do("ZRANGE", defaultSrvName, 0, -1))
	if err != nil {
		return "", err
	}

	for _, addr := range values {

		if addr == "" {
			log.Info("addr is a empty string:")
			continue
		}

		util := GetSrvLoadUtil(addr)
		tm, err := util.GetAddrUpdateTime()
		if err != nil {
			return "", errors.New("redis is invalid")
		}

		now := time.Now().UnixNano() / 1e6
		if now-tm > TIMEOUT {
			c.Do("ZREM", util.srvName, util.srvAddr)
			log.Info("addr is timeout: ", addr)
			continue
		}

		_, err = c.Do("ZINCRBY", defaultSrvName, 1, util.srvAddr)
		if err != nil {
			log.Error("serverloadutil zadd exec failed! err: ", err)
		}
		return util.srvAddr, nil
	}

	return "", errors.New("server not found")
}

// UpdateLoad 更新load
func (util *SrvLoadUtil) UpdateLoad(load int) error {
	now := time.Now().UnixNano() / 1e6
	err := util.SetAddrUpdateTime(now)
	if err != nil {
		return err
	}

	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err = c.Do("ZADD", util.srvName, load, util.srvAddr)
	if err != nil {
		return err
	}

	return nil
}

// SetAddrUpdateTime 设置更新时间
func (util *SrvLoadUtil) SetAddrUpdateTime(now int64) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("SET", util.key(), now)
	if err != nil {
		return err
	}

	return nil
}

// GetAddrUpdateTime 获得更新时间
func (util *SrvLoadUtil) GetAddrUpdateTime() (int64, error) {
	c := dbservice.GetCacheConn()
	defer c.Close()

	tm, err := redis.Int64(c.Do("GET", util.key()))
	if err != nil {
		return 0, err
	}

	return tm, nil
}

// IsValidAddr 是否抵制有效
func (util *SrvLoadUtil) IsValidAddr() error {
	tm, err := util.GetAddrUpdateTime()
	if err != nil {
		log.Error("IsValidAddr, err: ", err)
		return errors.New("redis is invalid")
	}

	now := time.Now().UnixNano() / 1e6
	if now-tm > TIMEOUT {
		log.Info("addr is timeout: ", util.srvAddr)
		return errors.New("server timeout")
	}

	return nil
}

// key 获取key
func (util *SrvLoadUtil) key() string {
	return fmt.Sprintf("%s:%s", util.srvName, util.srvAddr)
}
