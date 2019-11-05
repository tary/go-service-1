package dbservice

import (
	"encoding/json"
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

const (
	matchPrefix = "Match"
	teamPrefix  = "MatchTeam"
)

// MatchInfo 匹配信息
type MatchInfo struct {
	MatchType int32  //匹配类型
	MapID     uint32 //地图id
	MatchAddr string //匹配服务器地址
	ServerID  uint64 //匹配服务器的ID
}

// MatchUtil match相关
type MatchUtil struct {
	uid uint64
}

// GetMatchUtil 匹配信息工具类
func GetMatchUtil(uid uint64) *MatchUtil {
	return &MatchUtil{
		uid: uid,
	}
}

//key key
func (r *MatchUtil) key() string {
	return fmt.Sprintf("%s:%d", teamPrefix, r.uid)
}

// SaveMatchInfo 保存匹配信息
func (r *MatchUtil) SaveMatchInfo(info *MatchInfo) {
	if info == nil {
		log.Warn("SaveMatchInfo error info = ", info)
		return
	}

	d, e := json.Marshal(info)
	if e != nil {
		log.Warn("SaveMatchInfo error e = ", e)
	}
	if err := dbservice.CacheHSET(r.key(), r.uid, string(d)); err != nil {
		log.Error(err)
	}
}

// GetMatchInfo 获取匹配信息
func (r *MatchUtil) GetMatchInfo() *MatchInfo {
	exists, err := dbservice.CacheHExists(r.key(), r.uid)
	if err != nil || !exists {
		log.Error("GetMatchInfo CacheHExists failed: ", err)
		return nil
	}

	v, err := redis.String(dbservice.CacheHGET(r.key(), r.uid))
	if err != nil {
		log.Error("GetMatchInfo CacheHGET failed: ", err)
		return nil
	}

	var d *MatchInfo
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Warn("GetMatchInfo Failed to Unmarshal ", err)
		return nil
	}
	return d
}

// DelMatchInfo 删除匹配信息
func (r *MatchUtil) DelMatchInfo() {
	if err := dbservice.CacheHDEL(r.key(), r.uid); err != nil {
		log.Error("uid :", r.uid, " err:", err)
	}
}
