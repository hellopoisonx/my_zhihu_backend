package dao

import (
	"context"
	"errors"
	"fmt"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/log"
	"my_zhihu_backend/app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var l = log.L().With(zap.String("module", "user dao"))

type UserDAO struct {
	cfg config.ReadConfigFunc
	db  *gorm.DB
}

func NewUserDAO(cfg config.ReadConfigFunc, db *gorm.DB) *UserDAO {
	return &UserDAO{cfg: cfg, db: db}
}

// CreateNewUser 创建新用户
func (dao *UserDAO) CreateNewUser(ctx context.Context, user *model.User) app_error.AppError {
	err := gorm.G[model.User](dao.db).Create(ctx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return app_error.ErrUserAlreadyExists
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return nil
}

// GetById 通过 model.UserId 获取用户详情
func (dao *UserDAO) GetById(ctx context.Context, id model.UserId) (*model.User, app_error.AppError) {
	user, err := gorm.G[model.User](dao.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrUserNotExists
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return &user, nil
}

// GetByEmail 通过 email[string] 获取用户详情
func (dao *UserDAO) GetByEmail(ctx context.Context, email string) (*model.User, app_error.AppError) {
	user, err := gorm.G[model.User](dao.db).Where("email = ?", email).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrUserNotExists.WithError(err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return &user, nil
}

// ListUserByUsername 通过 username 搜索用户 返回用户详情列表
func (dao *UserDAO) ListUserByUsername(ctx context.Context, username string) ([]model.User, app_error.AppError) {
	users, err := gorm.G[model.User](dao.db).Where("username LIKE ?", "%"+username+"%").Find(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return users, nil
}

var ErrMysqlInvalidFields = app_error.NewInputError("invalid fields", app_error.ErrCodeInvalidParameters, nil)

// UpdateFields 更新用户信息
func (dao *UserDAO) UpdateFields(ctx context.Context, id model.UserId, fields map[string]any) app_error.AppError {
	if other, exists := fields["other"]; exists {
		otherMap, ok := other.(map[string]any)
		if !ok {
			return ErrMysqlInvalidFields
		}
		if len(otherMap) <= 0 {
			delete(fields, "other") // 防止other为空值时更新把原本json数据清空
		} else {
			expr := "json_set(ifnull(other, '{}')"
			args := make([]any, 0, len(otherMap)*2)
			for k, v := range otherMap {
				expr += ", ?, ?"
				args = append(args, "$."+k, v)
			}
			expr += ")"
			fields["other"] = gorm.Expr(expr, args...)
		}
	}
	if settings, exists := fields["settings"]; exists {
		otherMap, ok := settings.(map[string]any)
		if !ok {
			return ErrMysqlInvalidFields
		}
		if len(otherMap) <= 0 {
			delete(fields, "settings") // 防止settings为空值时更新把原本json数据清空
		} else {
			expr := "json_set(ifnull(settings, '{}')"
			args := make([]any, 0, len(otherMap)*2)
			for k, v := range otherMap {
				if b, ok := v.(bool); ok {
					expr += fmt.Sprintf(", ?, %t", b)
					args = append(args, "$."+k)
				} else {
					expr += ", ?, ?"
					args = append(args, "$."+k, v)
				}
			}
			expr += ")"
			fields["settings"] = gorm.Expr(expr, args...)
		}
	}
	result := dao.db.WithContext(ctx).Model(model.User{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return app_error.NewInternalError(app_error.ErrCodeMysql, result.Error)
	}
	if result.RowsAffected != 1 {
		return app_error.ErrUserNotExists
	}
	return nil
}

// DeleteUser 删除用户（软删除）
func (dao *UserDAO) DeleteUser(ctx context.Context, id model.UserId) app_error.AppError {
	rowsAffected, err := gorm.G[model.User](dao.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_error.ErrUserNotExists.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	if rowsAffected != 1 {
		return app_error.ErrUserNotExists
	}
	return nil
}

// FollowUser 添加关注关系
func (dao *UserDAO) FollowUser(ctx context.Context, followerID, followingID model.UserId) app_error.AppError {
	tx := dao.db.Begin()
	// 检查被关注用户是否存在
	_, err := gorm.G[model.User](tx).Where("id = ?", followingID).First(ctx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_error.ErrUserNotExists.WithError(err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	relation := model.UserFollowers{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	err = gorm.G[model.UserFollowers](tx).Create(ctx, &relation)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) { // 保证幂等性 重复关注不报错
			l.Warn("duplicate relationship", zap.Any("relation", relation))
			return nil
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	_, err = gorm.G[model.User](tx).
		Where("id = ?", followingID).
		Update(ctx, "follower_count", gorm.Expr("follower_count + ?", 1))
	if err != nil {
		tx.Rollback()
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	_, err = gorm.G[model.User](tx).
		Where("id = ?", followerID).
		Update(ctx, "following_count", gorm.Expr("following_count + ?", 1))
	if err != nil {
		tx.Rollback()
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	tx.Commit()
	return nil
}

// UnfollowUser 取消关注
func (dao *UserDAO) UnfollowUser(ctx context.Context, followerID, followingID model.UserId) app_error.AppError {
	tx := dao.db.Begin()
	relation := model.UserFollowers{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	rowsAffected, err := gorm.G[model.UserFollowers](tx).Where("follower_id = ? and following_id = ?", followerID, followingID).Delete(ctx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) { // 不报错
			l.Warn("duplicate relationship", zap.Any("relation", relation))
			return nil
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	if rowsAffected != 1 {
		tx.Rollback() // 不报错
		return nil
	}

	_, err = gorm.G[model.User](tx).
		Where("id = ?", followingID).
		Update(ctx, "follower_count", gorm.Expr("follower_count - ?", 1))
	if err != nil {
		tx.Rollback()
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	_, err = gorm.G[model.User](tx).
		Where("id = ?", followerID).
		Update(ctx, "following_count", gorm.Expr("following_count - ?", 1))
	if err != nil {
		tx.Rollback()
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	tx.Commit()
	return nil
}

// ListFollowers 获取用户的粉丝列表
func (dao *UserDAO) ListFollowers(ctx context.Context, followingID model.UserId) ([]model.UserId, app_error.AppError) {
	followers, err := gorm.G[model.UserFollowers](dao.db).Select("follower_id").Where("following_id = ?", followingID).Find(ctx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	results := make([]model.UserId, 0, len(followers))
	for _, follower := range followers {
		results = append(results, follower.FollowerID)
	}

	return results, nil
}

// ListFollowings 获取用户关注的用户列表
func (dao *UserDAO) ListFollowings(ctx context.Context, followerID model.UserId) ([]model.UserId, app_error.AppError) {
	followings, err := gorm.G[model.UserFollowers](dao.db).Select("following_id").Where("follower_id = ?", followerID).Find(ctx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	results := make([]model.UserId, 0, len(followings))
	for _, following := range followings {
		results = append(results, following.FollowingID)
	}

	return results, nil
}
