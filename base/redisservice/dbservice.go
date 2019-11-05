package dbservice

import (
	"fmt"
	"sync"
	"time"

	assert "github.com/aurelien-rainone/assertgo"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
)

var (
	// 仅初始化一次
	once       sync.Once
	onceServer sync.Once

	// redis连接池
	pool *redis.Pool

	// isDBValid DB是否正常
	isDBValid = atomic.NewBool(true)

	// 给服务器间同步状态使用的redis连接池
	poolForServer *redis.Pool
)

// GetPersistenceConn 获取一个持久化的redis连接
func GetPersistenceConn() redis.Conn {
	once.Do(initPool)
	return pool.Get()
}

// GetCacheConn 获取一个非持久化的redis连接
func GetCacheConn() redis.Conn {
	onceServer.Do(initServerPool)
	return poolForServer.Get()
}

// IsDBRedisValid DB redis是否可用
func IsDBRedisValid() bool {
	c := GetPersistenceConn()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		return false
	}

	return true
}

// IsServerRedisValid ForServer的redis是否可用
func IsServerRedisValid() bool {
	c := GetCacheConn()
	defer c.Close()

	if _, err := c.Do("PING"); err != nil {
		return false
	}

	return true
}

// initPool 初始化, 创建redis连接池
// 配置文件中需要配置:
// RedisServerAddr
// MaxIdle
// IdleTimeout
func initPool() {
	assert.True(pool == nil, "pool already inited")
	addr := viper.GetString("DB.Addr")
	index := viper.GetString("DB.Index")
	rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
	pwd := viper.GetString("DB.Password")

	maxIdle := viper.GetInt("DB.MaxIdle")
	idleTimeout := viper.GetInt("DB.IdleTimeout")
	maxActive := viper.GetInt("DB.MaxActive")
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		Wait:        true,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			if pwd != "" {
				return redis.DialURL(rawURL, redis.DialPassword(pwd))
			}
			return redis.DialURL(rawURL)
		},
	}

	go checkHealth()
}

// initServerPool 初始化给服务器同步信息用的redis
func initServerPool() {
	if poolForServer == nil {
		addr := viper.GetString("RedisForServer.Addr")
		index := viper.GetString("RedisForServer.Index")
		rawURL := fmt.Sprintf("redis://%s/%s", addr, index)
		pwd := viper.GetString("RedisForServer.Password")
		maxIdle := viper.GetInt("RedisForServer.MaxIdle")
		idleTimeout := viper.GetInt("RedisForServer.IdleTimeout")
		maxActive := viper.GetInt("RedisForServer.MaxActive")
		poolForServer = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			Wait:        true,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				if pwd != "" {
					return redis.DialURL(rawURL, redis.DialPassword(pwd))
				}
				return redis.DialURL(rawURL)
			},
		}
	}
}

// checkHealth 检查健康
func checkHealth() {
	t := viper.GetInt("Config.RedisHealthCheckTimer")
	if t == 0 {
		t = 1
	}

	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkDBValid()
		}
	}
}

// IsDBValid db是否有效
func IsDBValid() bool {
	return isDBValid.Load()
}

// checkDBValid 检查db有效性
func checkDBValid() {
	// DBValid DB是否正常
	DBValid := IsDBRedisValid()
	// SrvRedisValid Cache是否正常
	SrvRedisValid := IsServerRedisValid()

	isDBValid.Store(DBValid && SrvRedisValid)
}
