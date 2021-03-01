package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/gomodule/redigo/redis"
	"github.com/streadway/amqp"
)

// FailOnError raise to Sentry and fail on error
func FailOnError(err error, msg string) {
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("reported to Sentry: %s", err)
		log.Fatalf("%s: %s", msg, err)
	}
}

// SetupSentry run initial Sentry configuration
func SetupSentry() {
	sampleRate, err := strconv.ParseFloat(os.Getenv("SENTRY_SAMPLE_RATE"), 64)
	err = sentry.Init(sentry.ClientOptions{
		Dsn:        os.Getenv("SENTRY_DSN"),
		SampleRate: sampleRate,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

// NegAckAndCapture negate acknowledgement and capture error
func NegAckAndCapture(d amqp.Delivery, err error) {
	// Negatively acknowledge delivery and don't requeue plus break
	d.Nack(false, false)
	sentry.CaptureException(err)
	log.Printf("reported to Sentry: %s", err)
}

// NewRedisPool create Redis connection pool
func NewRedisPool() *redis.Pool {
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

// RedisSetEx helper for single SETEX until Redis 6.2 stable. Returns true if already exists
func RedisSetEx(rdb redis.Conn, key string, expireSeconds int) (bool, error) {
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
