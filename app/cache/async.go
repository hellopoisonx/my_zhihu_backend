package cache

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"time"
)

type Result[T any] struct {
	value chan *T
	error chan app_error.AppError
}

func NewResult[T any]() *Result[T] {
	return &Result[T]{
		value: make(chan *T, 1), // 必须添加缓冲区 否则会导致死锁
		error: make(chan app_error.AppError, 1),
	}
}

func (r *Result[T]) Value(ctx context.Context) (*T, app_error.AppError) {
	select {
	case <-ctx.Done():
		return nil, app_error.ErrTimeout
	case v := <-r.value:
		return v, nil
	case err := <-r.error:
		return nil, err
	}
}

func (r *Result[T]) Err(ctx context.Context) app_error.AppError {
	select {
	case <-ctx.Done():
		return app_error.ErrTimeout
	case err := <-r.error:
		return err
	}
}

type AsyncCacher[T any] struct {
	base Cacher[T]
}

func NewAsyncCacher[T any](base Cacher[T]) *AsyncCacher[T] {
	return &AsyncCacher[T]{
		base: base,
	}
}

func (a *AsyncCacher[T]) BaseCacher() Cacher[T] {
	return a.base
}

func (a *AsyncCacher[T]) PutChan(ctx context.Context, key string, value T) *Result[T] {
	result := NewResult[T]()
	go func() {
		result.error <- a.base.Put(ctx, key, value)
	}()
	return result
}

func (a *AsyncCacher[T]) GetChan(ctx context.Context, key string) *Result[T] {
	result := NewResult[T]()
	go func() {
		if value, err := a.base.Get(ctx, key); err != nil {
			result.error <- err
		} else {
			result.value <- value
		}
	}()
	return result
}

func (a *AsyncCacher[T]) RenewChan(ctx context.Context, key string, ttl time.Duration) *Result[T] {
	result := NewResult[T]()
	go func() {
		result.error <- a.base.Renew(ctx, key, ttl)
	}()
	return result
}

func (a *AsyncCacher[T]) InvalidateChan(ctx context.Context, key string) *Result[T] {
	result := NewResult[T]()
	go func() {
		if value, err := a.base.Invalidate(ctx, key); err != nil {
			result.error <- err
		} else {
			result.value <- value
		}
	}()
	return result
}
