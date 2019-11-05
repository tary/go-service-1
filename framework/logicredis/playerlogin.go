package dbservice

import (
	"encoding/json"
	"fmt"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// LoginData 登录临时数据
type LoginData struct {
	LoginTime int64  `json:"logintime"` // 登录时间
	GameAddr  string `json:"gameaddr"`  // 游戏地址
}

const (
	loginDataKey      = "LoginData"
	loginLockKey      = "LoginLockKey"
	lobbyAddr         = "LobbyAddr"
	playerLoginPrefix = "PlayerLoginPrefix"
	loginLockValue    = "V"
)

// PlayerLoginUtil player登录相关
type PlayerLoginUtil struct {
	uid uint64
}

// GetPlayerLoginUtil 登录工具类
func GetPlayerLoginUtil(uid uint64) *PlayerLoginUtil {
	return &PlayerLoginUtil{
		uid: uid,
	}
}

// SetLoginData 设置登录数据
func (util *PlayerLoginUtil) SetLoginData(data *LoginData) error {
	assert.True(data != nil, "LoginData is nil")

	d, e := json.Marshal(data)
	if e != nil {
		log.Error("SetLoginData err: ", e, "data: ", data)
	}

	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("SET", util.key(), d)

	return err
}

// GetLoginData 获取登录数据
func (util *PlayerLoginUtil) GetLoginData() *LoginData {
	c := dbservice.GetCacheConn()
	defer c.Close()

	reply, err := c.Do("GET", util.key())
	if err != nil {
		log.Error("GetLoginData err: ", err)
		return nil
	}

	v, err := redis.String(reply, nil)
	if err != nil {
		return nil
	}

	var d *LoginData
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Error("GetLoginData Failed to Unmarshal: ", err)
		return nil
	}

	return d
}

// DelLoginData 删除登录数据
func (util *PlayerLoginUtil) DelLoginData() error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("DEL", util.key())
	if err != nil {
		log.Error("GetLoginData err: ", err)
		return err
	}

	return nil
}

// Lock lock
func (util *PlayerLoginUtil) Lock() error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("SET", util.lockKey(), loginLockValue, "EX", 5, "NX")
	if err != nil {
		log.Error("Lock err: ", err)
		return err
	}

	return nil
}

// Unlock unlock
func (util *PlayerLoginUtil) Unlock() error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("DEL", util.lockKey())
	if err != nil {
		log.Error("Unlock err: ", err)
		return err
	}

	return nil
}

// key key
func (util *PlayerLoginUtil) key() string {
	return fmt.Sprintf("%s:%d", playerLoginPrefix, util.uid)
}

// lockkey lockkey
func (util *PlayerLoginUtil) lockKey() string {
	return fmt.Sprintf("%s:%d", loginLockKey, util.uid)
}
