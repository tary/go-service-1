package dbservice

import (
	"time"
)

// GlobalAtomicRun 使用Redis分布锁，跨进程，原子执行函数
func GlobalAtomicRun(lockName string, expireSec int64, bAsync bool, f func()) error {
	c := GetCacheConn()
	defer c.Close()

	body := func() error {
		for {
			_, err := c.Do("SET", lockName, 1, "EX", expireSec, "NX")
			// 获得锁
			if err == nil {
				break
			}

			// 等待锁
			t := time.NewTimer(1 * time.Millisecond)
			<-t.C
		}
		// 已获得锁
		f()
		// 释放锁
		c.Do("DEL", lockName)
		return nil
	}

	// 异步执行
	if bAsync {
		go body()
	} else {
		return body()
	}

	return nil
}
