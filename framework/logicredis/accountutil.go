package dbservice

import (
	"fmt"
	"math"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	dbservice "github.com/giant-tech/go-service/base/redisservice"
)

// AccountUtil 账号相关
type AccountUtil struct {
	uid uint64
}

const (
	// AccountPrefix 帐号表前缀
	AccountPrefix = "account"

	// AccountOpenID 用户名表前缀, 存储帐号和UID的对应关系
	AccountOpenID = "accountopenid"

	// UIDField UID字段
	UIDField = "uid"
)

// GetUID 通过username获取uid
func GetUID(user string) (uint64, error) {
	c := dbservice.GetPersistenceConn()
	defer c.Close()

	return redis.Uint64(c.Do("HGET", AccountOpenID+":"+user, UIDField))
}

// key key
func (util *AccountUtil) key() string {
	return fmt.Sprintf("%s:%d", AccountPrefix, util.uid)
}

// Account 获取帐号表工具类
func Account(uid uint64) *AccountUtil {
	acc := &AccountUtil{}
	acc.uid = uid
	return acc
}

// SetPassword SetPassword
func (util *AccountUtil) SetPassword(password string) error {
	c := dbservice.GetPersistenceConn()
	defer c.Close()

	_, err := c.Do("HSET", util.key(), "password", password)
	return err
}

// VerifyPassword 验证密码
func (util *AccountUtil) VerifyPassword(password string) bool {
	pwd, err := util.getPassword()
	if err != nil {
		log.Error("Get password failed", err)
		return false
	}

	if pwd != password {
		log.Error("Password not match")
		return false
	}

	return true
}

// SetUsername 保存用户名
func (util *AccountUtil) SetUsername(user string) error {
	c := dbservice.GetPersistenceConn()
	defer c.Close()

	if reply, err := redis.Int(c.Do("HSETNX", util.key(), "username", user)); err != nil {
		return err
	} else if reply == 0 {
		return fmt.Errorf("Account existed %s", user)
	}

	_, err := c.Do("HSET", AccountOpenID+":"+user, UIDField, util.uid)
	return err
}

// GetUsername 获得用户名
func (util *AccountUtil) GetUsername() (string, error) {
	c := dbservice.GetPersistenceConn()
	defer c.Close()
	return redis.String(c.Do("HGET", util.key(), "username"))
}

// SetGrade 设置帐号级别 1表示内部 2表示外部玩家
func (util *AccountUtil) SetGrade(grade uint32) error {
	c := dbservice.GetPersistenceConn()
	defer c.Close()

	_, err := c.Do("HSET", util.key(), "grade", grade)
	return err
}

// GetGrade 获取帐号级别
func (util *AccountUtil) GetGrade() (uint32, error) {

	c := dbservice.GetPersistenceConn()
	defer c.Close()
	v, err := c.Do("HGET", util.key(), "grade")
	if err != nil {
		return math.MaxUint32, err
	}
	if v == nil {
		return 0, nil
	}
	grade, err := redis.Uint64(v, nil)
	if err != nil {
		return math.MaxUint32, err
	}
	return uint32(grade), nil
}

// getPassword 方法根据用户名返回密码
func (util *AccountUtil) getPassword() (string, error) {
	c := dbservice.GetPersistenceConn()
	defer c.Close()

	return redis.String(c.Do("HGET", util.key(), "password"))
}
