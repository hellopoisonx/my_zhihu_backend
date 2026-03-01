package cache

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/util"
	"time"

	"github.com/redis/go-redis/v9"
)

type JsonCacher[T any] struct {
	prefix      string
	client      *redis.Client
	ttl         time.Duration
	fallback    Fallback[T]
	plainCacher *PlainCacher[string]
	bloomFilter *BloomFilter
}

func (cacher *JsonCacher[T]) BloomFilter() *BloomFilter {
	return cacher.bloomFilter
}

func NewJsonCacher[T any](client *redis.Client, ttl time.Duration, prefix string, fallback Fallback[T], filter *BloomFilter) *JsonCacher[T] {
	ttl += time.Duration(rand.IntN(5)) * time.Second // 增加过期时间的随机性 防止缓存雪崩
	return &JsonCacher[T]{
		client:      client,
		ttl:         ttl,
		prefix:      prefix,
		fallback:    fallback,
		bloomFilter: filter,
		plainCacher: NewPlainCacher[string](client, ttl, prefix, func(ctx context.Context, args ...any) (*string, app_error.AppError) {
			value, err := fallback(ctx, args...)
			if err != nil {
				return nil, err
			}
			if rawJson, err := json.Marshal(value); err != nil {
				return nil, app_error.ErrRedisCache.WithError(err)
			} else {
				return util.Ptr(string(rawJson)), nil
			}
		}, filter),
	}
}

func (cacher *JsonCacher[T]) Fallback() Fallback[T] {
	return cacher.fallback
}

func (cacher *JsonCacher[T]) TTL() time.Duration {
	return cacher.ttl
}

func (cacher *JsonCacher[T]) Prefix() string {
	return cacher.prefix
}

func (cacher *JsonCacher[T]) Renew(ctx context.Context, key string, ttl time.Duration) app_error.AppError {
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

func (cacher *JsonCacher[T]) Invalidate(ctx context.Context, key string) (*T, app_error.AppError) {
	if rawJson, err := cacher.client.GetDel(ctx, cacher.prefix+key).Bytes(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, app_error.ErrRedisCacheKeyNotExists.WithError(err)
		}
		return nil, app_error.ErrRedisCache.WithError(err)
	} else {
		value := new(T)
		if err := json.Unmarshal(rawJson, value); err != nil {
			return nil, app_error.ErrRedisCache.WithError(err)
		}
		return value, nil
	}
}

func (cacher *JsonCacher[T]) Get(ctx context.Context, key string, args ...any) (*T, app_error.AppError) {
	if rawJson, err := cacher.plainCacher.Get(ctx, key, args...); err != nil {
		return nil, err
	} else {
		value := new(T)
		if err := json.Unmarshal([]byte(*rawJson), value); err != nil {
			return nil, app_error.ErrRedisCache.WithError(err)
		}
		return value, nil
	}
}

func (cacher *JsonCacher[T]) Put(ctx context.Context, key string, value T) app_error.AppError {
	rawJson, err := json.Marshal(&value)
	if err != nil {
		return app_error.ErrRedisCache.WithError(err)
	}
	return cacher.plainCacher.Put(ctx, key, string(rawJson))
}
