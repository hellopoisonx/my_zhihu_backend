package repository

import (
	"my_zhihu_backend/app/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.ReadConfigFunc) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg().Redis.Addr,
	})
	return client
}
