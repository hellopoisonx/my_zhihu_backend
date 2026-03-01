package cache

import (
	"context"
	"errors"
	"my_zhihu_backend/app/app_error"

	"github.com/redis/go-redis/v9"
)

type BloomFilter struct {
	client *redis.Client
	name   string
}

func NewBloomFilter(name string, client *redis.Client) *BloomFilter {
	return &BloomFilter{
		client: client,
		name:   name,
	}
}

func (b *BloomFilter) Init(ctx context.Context) app_error.AppError {
	if err := b.client.BFReserve(ctx, b.name, 0.05, 1000000).Err(); err != nil {
		return app_error.ErrBloomFilter.WithError(err)
	}
	return nil
}

func (b *BloomFilter) Add(ctx context.Context, value any) app_error.AppError {
	if _, err := b.client.BFAdd(ctx, b.name, value).Result(); err != nil {
		return app_error.ErrBloomFilter.WithError(err)
	}
	return nil
}

func (b *BloomFilter) MAdd(ctx context.Context, values ...any) app_error.AppError {
	if _, err := b.client.BFMAdd(ctx, b.name, values...).Result(); err != nil {
		return app_error.ErrBloomFilter.WithError(err)
	}
	return nil
}

// AddChan 业务执行时在后台将数据异步存入布隆过滤器
func (b *BloomFilter) AddChan(ctx context.Context, value any) <-chan app_error.AppError {
	errC := make(chan app_error.AppError, 1)
	go func() {
		if _, err := b.client.BFAdd(ctx, b.name, value).Result(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				errC <- app_error.ErrTimeout.WithError(err)
			} else {
				errC <- app_error.ErrBloomFilter.WithError(err)
			}
		}
		errC <- nil
	}()
	return errC
}

// MAddChan 业务执行时在后台将数据异步存入布隆过滤器
func (b *BloomFilter) MAddChan(ctx context.Context, values ...any) <-chan app_error.AppError {
	errC := make(chan app_error.AppError, 1)
	go func() {
		if _, err := b.client.BFMAdd(ctx, b.name, values...).Result(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				errC <- app_error.ErrTimeout.WithError(err)
			} else {
				errC <- app_error.ErrBloomFilter.WithError(err)
			}
		}
		errC <- nil
	}()
	return errC
}

func (b *BloomFilter) Exist(ctx context.Context, value any) (bool, app_error.AppError) {
	if exists, err := b.client.BFExists(ctx, b.name, value).Result(); err != nil {
		return false, app_error.ErrBloomFilter.WithError(err)
	} else {
		return exists, nil
	}
}
