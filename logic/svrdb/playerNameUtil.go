package svrdb

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// 玩家在线情况下各种临时信息的存储, 不需要持久化
// 以玩家名作为key

const (
	playerNamePrefix     = "PlayerName"
	playerServerIDPrefix = "ServerID"
)

// PlayerNameUtil PlayerNameUtil
type PlayerNameUtil struct {
	name string
}

// GetPlayerNameUtil 玩家名字相关
func GetPlayerNameUtil(name string) *PlayerNameUtil {
	return &PlayerNameUtil{
		name: name,
	}
}

// key key
func (u *PlayerNameUtil) key() string {
	return fmt.Sprintf("%s:%s", playerNamePrefix, u.name)
}

// GetPlayerServerID 获取玩家serverID
func (u *PlayerNameUtil) GetPlayerServerID(serverType int32) uint64 {

	c := dbservice.GetCacheConn()
	defer c.Close()

	reply, err := c.Do("HGET", u.key(), fmt.Sprintf("%s:%d", playerServerIDPrefix, serverType))
	if err != nil {
		return 0
	}

	var state uint64
	state, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0
	}

	return state
}

// SetPlayerServiceID 设置玩家serverID
func (u *PlayerNameUtil) SetPlayerServiceID(serverType int32, serverID uint64) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("HSET", u.key(), fmt.Sprintf("%s:%d", playerServerIDPrefix, serverType), serverID)

	return err
}

// DelPlayerServiceID 删除玩家serverID
func (u *PlayerNameUtil) DelPlayerServiceID(serverType int32) error {
	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("HDEL", u.key(), fmt.Sprintf("%s:%d", playerServerIDPrefix, serverType))

	return err
}
