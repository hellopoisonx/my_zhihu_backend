package cache

import (
	"cmp"
	"context"
	"errors"
	"math/rand/v2"
	"my_zhihu_backend/app/app_error"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/redis/go-redis/v9"
)

type PlainCacher[T cmp.Ordered] struct {
	prefix   string
	client   *redis.Client
	ttl      time.Duration
	fallback Fallback[T]

	bloomFilter *BloomFilter        // 解决缓存穿透 拦截恶意请求
	single      *singleflight.Group // 使用 singleflight 包实现同一个key同一时间只能有一个 fallback 在执行 防止缓存击穿
}

func (cacher *PlainCacher[T]) BloomFilter() *BloomFilter {
	return cacher.bloomFilter
}

func NewPlainCacher[T cmp.Ordered](client *redis.Client, ttl time.Duration, prefix string, fallback Fallback[T], filter *BloomFilter) *PlainCacher[T] {
	ttl += time.Duration(rand.IntN(5)) * time.Second // 增加过期时间的随机性 防止缓存雪崩
	return &PlainCacher[T]{
		client:      client,
		ttl:         ttl,
		prefix:      prefix,
		fallback:    fallback,
		bloomFilter: filter,
		single:      new(singleflight.Group),
	}
}

func (cacher *PlainCacher[T]) Fallback() Fallback[T] {
	return cacher.fallback
}

func (cacher *PlainCacher[T]) TTL() time.Duration {
	return cacher.ttl
}

func (cacher *PlainCacher[T]) Prefix() string {
	return cacher.prefix
}

func (cacher *PlainCacher[T]) Renew(ctx context.Context, key string, ttl time.Duration) app_error.AppError {
	ttl += time.Duration(rand.IntN(5)) * time.Second // 增加过期时间的随机性 防止缓存雪崩
	if ok, err := cacher.client.Expire(ctx, cacher.prefix+key, ttl).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return app_error.ErrRedisCacheKeyNotExists.WithError(err)
		}
		return app_error.ErrRedisCache.WithError(err)
	} else if !ok {
		return app_error.ErrRedisCacheKeyNotExists
	}
	return nil
}

func (cacher *PlainCacher[T]) Invalidate(ctx context.Context, key string) (*T, app_error.AppError) {
	value := new(T)
	if err := cacher.client.GetDel(ctx, cacher.prefix+key).Scan(value); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, app_error.ErrRedisCacheKeyNotExists.WithError(err)
		}
		return nil, app_error.ErrRedisCache.WithError(err)
	}
	return value, nil
}

func (cacher *PlainCacher[T]) Put(ctx context.Context, key string, value T) app_error.AppError {
	if err := cacher.bloomFilter.Add(ctx, cacher.prefix+key); err != nil {
		return app_error.ErrRedisCache.WithError(err)
	}
	if err := cacher.client.Set(ctx, cacher.prefix+key, value, cacher.ttl).Err(); err != nil {
		return app_error.ErrRedisCache.WithError(err)
	}
	return nil
}

func (cacher *PlainCacher[T]) Get(ctx context.Context, key string, args ...any) (*T, app_error.AppError) {
	value := new(T)
	if exists, err := cacher.bloomFilter.Exist(ctx, cacher.prefix+key); err != nil {
		return nil, app_error.ErrRedisCache.WithError(err)
	} else if !exists {
		return nil, app_error.ErrRedisCacheKeyNotExists
	}

	if err := cacher.client.Get(ctx, cacher.prefix+key).Scan(value); err != nil {
		if errors.Is(err, redis.Nil) {
			res, err, _ := cacher.single.Do(key, func() (interface{}, error) {
				return cacher.fallback(ctx, args...)
			})
			if err != nil {
				return nil, err.(app_error.AppError)
			}
			if err := cacher.Put(ctx, key, *res.(*T)); err != nil {
				return nil, err
			}
			return res.(*T), nil
		}
		return nil, app_error.ErrRedisCache.WithError(err)
	}
	return value, nil
}
