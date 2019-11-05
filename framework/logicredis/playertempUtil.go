package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// 玩家在线情况下各种临时信息的存储, 不需要持久化

const (
	playerTempPrefix     = "PlayerTemp"
	playerGroupPrefix    = "PlayerGroup"
	playerServerIDPrefix = "ServerID"
)

// PlayerTempUtil player temp相关
type PlayerTempUtil struct {
	groupType uint32
	uid       uint64
}

// GetPlayerTempUtil 获取工具类
func GetPlayerTempUtil(uid uint64) *PlayerTempUtil {
	return &PlayerTempUtil{
		uid: uid,
	}
}

// key key
func (u *PlayerTempUtil) key() string {
	return fmt.Sprintf("%s:%d", playerTempPrefix, u.uid)
}

// groupKey groupKey
func (u *PlayerTempUtil) groupKey() string {
	return fmt.Sprintf("%s:%d", playerGroupPrefix, u.uid)
}

// GetPlayerGroupID 获取玩家所在group id
func (u *PlayerTempUtil) GetPlayerGroupID(groupType uint32) uint64 {
	reply, err := dbservice.CacheHGET(u.groupKey(), groupType)
	if err != nil {
		return 0
	}
	var groupID uint64
	groupID, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0
	}
	return groupID
}

// SetPlayerGroupID 设置玩家所在group id
func (u *PlayerTempUtil) SetPlayerGroupID(groupType uint32, groupID uint64) error {
	return dbservice.CacheHSET(u.groupKey(), groupType, groupID)
}

// DelPlayerGroupID 删除玩家所在group id
func (u *PlayerTempUtil) DelPlayerGroupID(groupType uint32) error {
	return dbservice.CacheHDEL(u.groupKey(), groupType)
}

// GetGameState 获取玩家游戏状态
func (u *PlayerTempUtil) GetGameState() uint64 {
	reply, err := dbservice.CacheHGET(u.key(), "gamestate")
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

// SetGameState 设置玩家游戏状态
func (u *PlayerTempUtil) SetGameState(state uint64) error {
	return dbservice.CacheHSET(u.key(), "gamestate", state)
}

// GetEnterGameTime 获取玩家进入游戏时间
func (u *PlayerTempUtil) GetEnterGameTime() uint64 {
	reply, err := dbservice.CacheHGET(u.key(), "entertime")
	if err != nil {
		return 0
	}
	var entertime uint64
	entertime, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0
	}
	return entertime
}

// SetEnterGameTime 设置玩家进入游戏时间
func (u *PlayerTempUtil) SetEnterGameTime(timestamp uint64) error {
	return dbservice.CacheHSET(u.key(), "entertime", timestamp)
}

// GetPlayerServerID 获取玩家serverID
func (u *PlayerTempUtil) GetPlayerServerID(serverType int32) uint64 {

	reply, err := dbservice.CacheHGET(u.key(), fmt.Sprintf("%s:%d", playerServerIDPrefix, serverType))
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

// SetPlayerServerID 设置玩家serverID
func (u *PlayerTempUtil) SetPlayerServerID(serverType int32, serverID uint64) error {
	return dbservice.CacheHSET(u.key(), fmt.Sprintf("%s:%d", playerServerIDPrefix, serverType), serverID)
}

// DelPlayerServerID 删除玩家serverID
func (u *PlayerTempUtil) DelPlayerServerID(serverType int32) error {
	return dbservice.CacheHDEL(u.key(), fmt.Sprintf("%s:%d", playerServerIDPrefix, serverType))
}
