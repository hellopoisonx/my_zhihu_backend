package cache

import (
	"context"
	"testing"
	"time"

	"my_zhihu_backend/app/app_error"

	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestJsonCacher_WithContainer(t *testing.T) {
	ctx := context.TODO()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // TODO: 使用 testcontainer 报错 目前回退本地redis服务
	})

	fallbackCalled := 0
	expectedUser := &TestUser{ID: 1, Name: "Tester"}

	bloom := NewBloomFilter("test", client)

	cacher := NewJsonCacher[TestUser](client, time.Minute, "test:", func(ctx context.Context) (*TestUser, app_error.AppError) {
		fallbackCalled++
		return expectedUser, nil
	}, bloom)

	err := bloom.Add(ctx, "test:user_1")
	assert.Nil(t, err)

	t.Run("Test_Get_And_Fallback", func(t *testing.T) {
		// 第一次获取：缓存不存在，触发 Fallback
		val, err := cacher.Get(ctx, "user_1")
		assert.Nil(t, err)
		assert.Equal(t, "Tester", val.Name)
		assert.Equal(t, 1, fallbackCalled)

		// 第二次获取：缓存已存在，不触发 Fallback
		val2, _ := cacher.Get(ctx, "user_1")
		assert.Equal(t, val, val2)
		assert.Equal(t, 1, fallbackCalled)
	})

	t.Run("Test_Invalidate", func(t *testing.T) {
		_, err := cacher.Invalidate(ctx, "user_1")

		assert.NoError(t, err)
	})
}
