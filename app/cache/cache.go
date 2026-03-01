package cache

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"time"
)

type Fallback[T any] func(ctx context.Context, args ...any) (*T, app_error.AppError) // 查询缓存不存在或者无效的时候执行 Fallback 回调并且将结果重新写入缓存

type Cacher[T any] interface {
	Put(ctx context.Context, key string, value T) app_error.AppError
	Get(ctx context.Context, key string, args ...any) (*T, app_error.AppError)
	Renew(ctx context.Context, key string, ttl time.Duration) app_error.AppError
	Invalidate(ctx context.Context, key string) (*T, app_error.AppError)

	Fallback() Fallback[T]
	TTL() time.Duration
	Prefix() string
	BloomFilter() *BloomFilter
}

type AsyncCacherDecorator[T any] interface {
	BaseCacher() Cacher[T]

	PutChan(ctx context.Context, key string, value T) *Result[T]
	GetChan(ctx context.Context, key string, args ...any) *Result[T]
	RenewChan(ctx context.Context, key string, ttl time.Duration) *Result[T]
	InvalidateChan(ctx context.Context, key string) *Result[T]
}
