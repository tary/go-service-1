package dbservice

import (
	"errors"
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// EntitySrvInfo 服务器信息
type EntitySrvInfo struct {
	SrvID   uint64
	GroupID uint64
}

// EntitySrvUtil 表工具类
// 存储entity分布在哪些服务器上的信息
// id为entityid
type EntitySrvUtil struct {
	id uint64
}

const (
	entitySrvInfoPrefix = "entitysrvinfo"
	fieldType           = "entitytype"
	fieldDBID           = "dbid"
	fieldExistedCnt     = "existedcnt"
)

// GetEntitySrvUtil 获得Entity工具类
func GetEntitySrvUtil(eid uint64) *EntitySrvUtil {
	enUtil := &EntitySrvUtil{}
	enUtil.id = eid
	return enUtil
}

// IsExist 这个ID号的Entity是否存在
func (util *EntitySrvUtil) IsExist() bool {
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

// RegSrvID 注册类型和服务信息
func (util *EntitySrvUtil) RegSrvID(srvType uint8, srvID uint64, groupID uint64, entityType string) error {
	//log.Debug("Reg service key: ", util.key(), ", serviceType: ", srvType, ", serviceID: ", srvID)

	c := dbservice.GetCacheConn()
	defer c.Close()

	_, err := c.Do("HMSET", util.key(), srvType, util.joinSrvInfo(srvID, groupID), fieldType, entityType)
	if err != nil {
		log.Error(err)
	}
	_, err = c.Do("HINCRBY", util.key(), fieldExistedCnt, 1)
	return err
}

// UnRegSrvID 删除注册信息
func (util *EntitySrvUtil) UnRegSrvID(srvType uint8, srvID uint64, groupID uint64) error {
	//log.Debug("UnReg service key: ", util.key(), ", serviceType: ", srvType, ", serviceID: ", srvID)

	c := dbservice.GetCacheConn()
	defer c.Close()

	reply, err := c.Do("HINCRBY", util.key(), fieldExistedCnt, -1)
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

	_srvID, _groupID, err := util.GetSrvInfo(srvType)
	if err != nil {
		return err
	}

	if _srvID == srvID && _groupID == groupID {
		_, err = c.Do("HDEL", util.key(), srvType)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

// GetEntityInfo 获取entityType
func (util *EntitySrvUtil) GetEntityInfo() (string, uint64, error) {

	c := dbservice.GetCacheConn()
	defer c.Close()

	ret, err := c.Do("HMGET", util.key(), fieldType, fieldDBID)
	if err != nil {
		return "", 0, err
	}

	retValue, err := redis.Values(ret, nil)
	if err != nil {
		return "", 0, err
	}

	if len(retValue) != 2 {
		return "", 0, errors.New("wrong")
	}

	entityType, err := redis.String(retValue[0], nil)
	if err != nil {
		return "", 0, err
	}

	dbID, err := redis.Uint64(retValue[1], nil)
	if err != nil {
		return "", 0, err
	}

	return entityType, dbID, nil
}

// GetSrvInfo 获取特定服务类型的 服务ID以及 GroupID
func (util *EntitySrvUtil) GetSrvInfo(srvType uint8) (srvID uint64, groupID uint64, err error) {

	c := dbservice.GetCacheConn()
	defer c.Close()

	ret, err := c.Do("HGET", util.key(), srvType)
	if err != nil {
		return
	}

	retStr, err := redis.String(ret, nil)
	if err != nil {
		return
	}

	srvID, groupID = util.splitSrvInfo(retStr)

	return
}

//GetGroupID 获取groupID
func (util *EntitySrvUtil) GetGroupID(srvType uint8) (groupID uint64) {
	c := dbservice.GetCacheConn()
	defer c.Close()

	ret, err := c.Do("HGET", util.key(), srvType)
	if err != nil {
		return
	}

	retStr, err := redis.String(ret, nil)
	if err != nil {
		return
	}

	_, groupID = util.splitSrvInfo(retStr)

	return
}

// GetSrvIDs 获取Entity的分布式信息
func (util *EntitySrvUtil) GetSrvIDs() (map[uint8]*EntitySrvInfo, error) {

	c := dbservice.GetCacheConn()
	defer c.Close()

	reply, err := c.Do("HGETALL", util.key())
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, nil
	}

	values, err := redis.Values(reply, nil)
	if err != nil {
		return nil, err
	}

	result := make(map[uint8]*EntitySrvInfo)

	for i := 0; i < len(values); i += 2 {
		srvType, err := redis.Uint64(values[i], nil)
		if err != nil {
			continue
		}

		s, _ := redis.String(values[i+1], nil)

		srvID, GroupID := util.splitSrvInfo(s)

		srvInfo := EntitySrvInfo{srvID, GroupID}

		result[uint8(srvType)] = &srvInfo
	}
	return result, nil
}

// GetCellInfo 获取包含Cell的srvID 和 groupID
func (util *EntitySrvUtil) GetCellInfo() (uint64, uint64, error) {

	srvInfos, err := util.GetSrvIDs()
	if err != nil {
		return 0, 0, err
	}

	for _, info := range srvInfos {
		if info.GroupID != 0 {
			return info.SrvID, info.GroupID, nil
		}
	}

	return 0, 0, nil
}

// joinSrvInfo 加入服务器信息
func (util *EntitySrvUtil) joinSrvInfo(t uint64, id uint64) string {
	return fmt.Sprintf("%d:%d", t, id)
}

// splitSrvInfo 拆分服务器信息
func (util *EntitySrvUtil) splitSrvInfo(s string) (srvID uint64, groupID uint64) {
	fmt.Sscanf(s, "%d:%d", &srvID, &groupID)
	return
}

// key key
func (util *EntitySrvUtil) key() string {
	return fmt.Sprintf("%s:%d", entitySrvInfoPrefix, util.id)
}
