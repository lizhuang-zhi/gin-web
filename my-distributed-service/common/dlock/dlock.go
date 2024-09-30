package dlock

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var redsyncInstance *redsync.Redsync

func InitDistruibutedLock(client *redis.Client) {
	// 创建一个 Redis 连接池
	pool := goredis.NewPool(client)
	// 使用连接池创建 redsync 实例
	redsyncInstance = redsync.New(pool)
}

func ObtainLock(key string, expiry int) (*redsync.Mutex, error) {
	mutex := redsyncInstance.NewMutex(key, redsync.WithExpiry(time.Duration(expiry)*time.Second))
	if err := mutex.Lock(); err != nil {
		return nil, err
	}

	return mutex, nil
}

func ReleaseLock(mutex *redsync.Mutex) error {
	_, err := mutex.Unlock()
	return err
}
