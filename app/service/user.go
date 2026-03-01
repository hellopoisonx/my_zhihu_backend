package service

import (
	"context"
	"fmt"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/cache"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/dao"
	"my_zhihu_backend/app/log"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/request"
	"my_zhihu_backend/app/util"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var l = log.L().With(zap.String("module", "user service"))

func NewUserService(db *gorm.DB, client *redis.Client) *UserService {
	cfg := config.C
	userDAO := dao.NewUserDAO(cfg, db)
	u := new(util.Util)
	bloomFilter := cache.NewBloomFilter("user-filter", client)
	infoCacher := cache.NewJsonCacher(client, 24*time.Hour, cfg().Prefix.UserInfoPrefix, func(ctx context.Context, args ...any) (*model.User, app_error.AppError) {
		return userDAO.GetById(ctx, args[0].(model.UserId))
	}, bloomFilter)
	return &UserService{
		dao:         userDAO,
		infoCacher:  infoCacher,
		bloomFilter: bloomFilter,
		cfg:         cfg,
		util:        u,
	}
}

type UserService struct {
	dao         *dao.UserDAO
	cfg         config.ReadConfigFunc
	util        *util.Util
	infoCacher  cache.Cacher[model.User]
	bloomFilter *cache.BloomFilter
}

func (service *UserService) store(_ context.Context, user model.User) {
	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		key := fmt.Sprintf("%d", user.Id)
		user.HPassword = ""
		res := cache.NewAsyncCacher(service.infoCacher).PutChan(timeout, key, user)
		if err := res.Err(timeout); err != nil {
			l.Error("failed to cache user info", err.ErrorField()...)
		} else {
			l.Info("user info catch", zap.Any("user_id", user.Id))
		}
	}()
}

func (service *UserService) CreateNewUser(ctx context.Context, req *request.CreateNewUserRequest) (*model.User, app_error.AppError) {
	if req.Region == "" {
		req.Region = "unknown"
	}
	if req.Other.Icon == nil {
		req.Other.Icon = util.Ptr("https://nothing.here")
	}
	if req.Other.Introduction == nil {
		req.Other.Introduction = util.Ptr("nothing here")
	}

	if req.Settings.HidePrivacy == nil {
		req.Settings.HidePrivacy = util.Ptr(false)
	}

	if hPasswd, err := service.util.EncryptPassword(req.Password); err != nil {
		return nil, app_error.NewInternalError(app_error.ErrCodeEncryption, err)
	} else {
		user := model.User{
			Id:             model.UserId(service.util.GenerateSnowflakeID()),
			Username:       req.Username,
			HPassword:      string(hPasswd),
			Email:          req.Email,
			FollowerCount:  0,
			FollowingCount: 0,
			Gender:         &req.Gender,
			Region:         req.Region,
			Settings:       model.UserSettings{HidePrivacy: *req.Settings.HidePrivacy},
			Other:          model.UserOtherInfo{Introduction: *req.Other.Introduction, Icon: *req.Other.Icon},
		}
		service.store(ctx, user)
		err := service.dao.CreateNewUser(ctx, &user)
		if err != nil {
			return nil, err
		}
		return service.dao.GetByEmail(ctx, req.Email)
	}
}

// DeleteUser 删除用户
func (service *UserService) DeleteUser(ctx context.Context, id int64) app_error.AppError {
	return service.dao.DeleteUser(ctx, model.UserId(id))
}

// GetUser 获取用户信息
func (service *UserService) GetUser(ctx context.Context, id int64) (*model.User, app_error.AppError) {
	return service.infoCacher.Get(ctx, fmt.Sprintf("%d", id))
}

// UpdateUser 更新用户信息
func (service *UserService) UpdateUser(ctx context.Context, id int64, req *request.UpdateUserRequest) (*model.User, app_error.AppError) {
	fields := make(map[string]interface{})

	if req.Username != "" {
		fields["username"] = req.Username
	}

	if req.Region != "" {
		fields["region"] = req.Region
	}

	if req.Gender != nil {
		fields["gender"] = *req.Gender
	}

	if req.Settings != nil {
		fields["settings"] = make(map[string]any)
		if req.Settings.HidePrivacy != nil {
			fields["settings"].(map[string]any)["hide_privacy"] = *req.Settings.HidePrivacy
		}
	}

	if req.Other != nil {
		fields["other"] = make(map[string]any)
		if req.Other.Introduction != nil {
			fields["other"].(map[string]any)["introduction"] = *req.Other.Introduction
		}
		if req.Other.Icon != nil {
			fields["other"].(map[string]any)["icon"] = *req.Other.Icon
		}
	}

	if req.Password != "" {
		hPasswd, err := service.util.EncryptPassword(req.Password)
		if err != nil {
			return nil, app_error.NewInternalError(app_error.ErrCodeEncryption, err)
		}
		fields["h_password"] = string(hPasswd)
	}

	if len(fields) == 0 {
		return nil, nil // 没有要更新的字段
	}

	user, err := service.dao.UpdateFields(ctx, model.UserId(id), fields)
	if err != nil {
		return nil, err
	}
	service.store(ctx, *user)
	return user, err
}

// SearchUserByUsername 根据用户名搜索用户
func (service *UserService) SearchUserByUsername(ctx context.Context, username string) ([]int64, app_error.AppError) {
	users, err := service.dao.ListUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(users))
	for _, user := range users {
		ids = append(ids, int64(user.Id))
	}
	return ids, nil
}

// FollowUser 关注用户
func (service *UserService) FollowUser(ctx context.Context, followerID, followingID int64) app_error.AppError {
	return service.dao.FollowUser(ctx, model.UserId(followerID), model.UserId(followingID))
}

// UnfollowUser 取消关注用户
func (service *UserService) UnfollowUser(ctx context.Context, followerID, followingID int64) app_error.AppError {
	return service.dao.UnfollowUser(ctx, model.UserId(followerID), model.UserId(followingID))
}

// GetFollowers 获取粉丝列表
func (service *UserService) GetFollowers(ctx context.Context, id model.UserId) ([]model.UserId, app_error.AppError) {
	return service.dao.ListFollowers(ctx, id)
}

// GetFollowings 获取关注的用户列表
func (service *UserService) GetFollowings(ctx context.Context, id model.UserId) ([]model.UserId, app_error.AppError) {
	return service.dao.ListFollowings(ctx, id)
}
