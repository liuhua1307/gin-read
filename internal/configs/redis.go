package configs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	_  Cache = &Redis{}
	rs *redis.Client
)

type Redis struct {
	rs *redis.Client
}

func init() {
	rs = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func NewRedis(rs *redis.Client) *Redis {
	return &Redis{rs: rs}
}

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	TryLock(key string, value string, expiration time.Duration) (bool, error)
	Unlock(key string) error
}

func NewCache(rs *redis.Client) Cache {
	return NewRedis(rs)
}
func (r *Redis) Get(key string) (string, error) {
	return r.rs.Get(context.Background(), key).Result()
}

func (r *Redis) Set(key string, value string) error {
	return r.rs.Set(context.Background(), key, value, 0).Err()
}

func (r *Redis) TryLock(key string, value string, expiration time.Duration) (bool, error) {
	return r.rs.SetNX(context.Background(), key, value, expiration).Result()
}

func (r *Redis) Unlock(key string) error {
	return r.rs.Del(context.Background(), key).Err()
}
