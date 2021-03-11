package redisutils

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

// NewPool create Redis connection pool
func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), // host:port
				redis.DialPassword(os.Getenv("REDIS_PASSWORD")),
			)
		},
	}
}

// SetNxEx helper for single SETEX until Redis 6.2 stable. Returns true if already exists
func SetNxEx(rdb redis.Conn, key string, expireSeconds int) (bool, error) {
	// TODO replace with 'SET key "true" EX time NX GET' when Redis v6.2 stable
	// Check if stored in Redis
	val, err := rdb.Do("SETNX", key, true)
	if err != nil {
		return false, err
	}
	if val != int64(1) {
		return true, nil
	}

	// Set key expiry
	_, err = rdb.Do("EXPIRE", key, 8*time.Hour.Seconds())
	if err != nil {
		return false, err
	}
	return false, nil
}

// Delete delete single key from redis connection
func Delete(rdb redis.Conn, key string) error {
	_, err := rdb.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}
