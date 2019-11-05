package dbservice

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

const (
	// groupPrefix prefix
	groupPrefix = "Group"
	// CenterIDStr str
	CenterIDStr = "ServiceID"
)

// GroupUtil 组相关
type GroupUtil struct {
	groupID uint64
}

// GetGroupUtil Group工具类(房间，队伍等)
func GetGroupUtil(groupID uint64) *GroupUtil {
	return &GroupUtil{
		groupID: groupID,
	}
}

// key key
func (r *GroupUtil) key() string {
	return fmt.Sprintf("%s:%d", groupPrefix, r.groupID)
}

// DelCenterSrvID 删除房间CenterID
func (r *GroupUtil) DelCenterSrvID() {
	dbservice.CacheHDEL(r.key(), CenterIDStr)
}

// SetCenterSrvID 设置房间所在的CenterID
func (r *GroupUtil) SetCenterSrvID(srvID uint64) error {
	return dbservice.CacheHSET(r.key(), CenterIDStr, srvID)
}

// GetCenterSrvID 获取房间所在的CenterID
func (r *GroupUtil) GetCenterSrvID() (uint64, error) {
	return redis.Uint64(dbservice.CacheHGET(r.key(), CenterIDStr))
}
