package dao

import (
	"context"
	"errors"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthDAO struct {
	client *redis.Client
	cfg    config.ReadConfigFunc
}

func NewAuthDAO(client *redis.Client, cfg config.ReadConfigFunc) *AuthDAO {
	return &AuthDAO{client: client, cfg: cfg}
}

func (dao *AuthDAO) SaveRefreshToken(ctx context.Context, id model.UserId, refreshToken string, exp time.Duration) app_error.AppError {
	if err := dao.client.Set(ctx, dao.cfg().Prefix.RefreshToken+strconv.Itoa(int(id)), refreshToken, exp).Err(); err != nil {
		return app_error.NewInternalError(app_error.ErrCodeRedis, err)
	}
	return nil
}

func (dao *AuthDAO) GetRefreshToken(ctx context.Context, id model.UserId) (string, app_error.AppError) {
	if token, err := dao.client.Get(ctx, dao.cfg().Prefix.RefreshToken+strconv.Itoa(int(id))).Result(); err != nil {
		return "", app_error.NewInternalError(app_error.ErrCodeRedis, err)
	} else {
		return token, nil
	}
}

func (dao *AuthDAO) DeleteRefreshToken(ctx context.Context, id model.UserId) app_error.AppError {
	if err := dao.client.Del(ctx, dao.cfg().Prefix.RefreshToken+strconv.Itoa(int(id))).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return app_error.NewInternalError(app_error.ErrCodeRedis, err)
	}
	return nil
}
