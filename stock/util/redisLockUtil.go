package util

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	retryDelay = 100 * time.Millisecond
)

type RedisLockClient struct {
	RedisClient *redis.Client
	Redsync     *redsync.Redsync
}

func (c *RedisLockClient) AcquireLock(key string, timeout time.Duration) (*redsync.Mutex, error) {
	mutex := c.Redsync.NewMutex(
		key,
		redsync.WithRetryDelay(retryDelay),
		redsync.WithTries(int(timeout/retryDelay)),
		// 25s后自动释放锁
		redsync.WithExpiry(25*time.Second),
	)
	return mutex, mutex.Lock()
}
