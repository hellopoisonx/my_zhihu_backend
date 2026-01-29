package mysql

import (
	"app/model"
	"app/proxy"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserDAO struct {
	*Mysql
	options map[string]any
	uuid    *string
}

func NewUserDAO(mysql *Mysql) proxy.UserDAOProxy {
	return &UserDAO{
		options: make(map[string]any),
		Mysql:   mysql,
		uuid:    nil,
	}
}

func (u *UserDAO) UUID(uuid string) proxy.UserDAOProxy {
	u.uuid = &uuid
	return u
}

func (u *UserDAO) Username(username string) proxy.UserDAOProxy {
	u.options["username"] = username
	return u
}

// Password password是已经经算法处理的加盐哈希值
// TODO: 增加检查机制
func (u *UserDAO) Password(password string) proxy.UserDAOProxy {
	u.options["password"] = password
	return u
}

func (u *UserDAO) Email(email string) proxy.UserDAOProxy {
	u.options["email"] = email
	return u
}

func (u *UserDAO) Gender(gender model.UserGender) proxy.UserDAOProxy {
	if gender < 0 || gender > 2 {
		u.log.Error("undefined gender", zap.Int("gender", int(gender)))
		u.options["gender"] = model.UserGenderSecret
	} else {
		u.options["gender"] = gender
	}
	return u
}

func (u *UserDAO) Region(region string) proxy.UserDAOProxy {
	u.options["region"] = region
	return u
}

func (u *UserDAO) Introduction(introduction string) proxy.UserDAOProxy {
	if _, exists := u.options["other"]; !exists {
		u.options["other"] = map[string]any{}
	}
	u.options["other"].(map[string]any)["introduction"] = introduction
	return u
}

func (u *UserDAO) Icon(icon string) proxy.UserDAOProxy {
	if _, exists := u.options["other"]; !exists {
		u.options["other"] = map[string]any{}
	}
	u.options["other"].(map[string]any)["icon"] = icon
	return u
}

func (u *UserDAO) Create(ctx context.Context) error {
	if u.uuid == nil {
		return errors.New("uuid is nil")
	}
	if _, exists := u.options["other"]; !exists {
		u.options["other"] = map[string]any{}
	}
	if _, exists := u.options["other"].(map[string]any)["introduction"]; !exists {
		u.options["other"].(map[string]any)["introduction"] = "nothing"
	}
	if _, exists := u.options["other"].(map[string]any)["icon"]; !exists {
		u.options["other"].(map[string]any)["icon"] = "" // TODO 默认图像url
	}
	return gorm.G[model.UserModel](u.db).Create(ctx, &model.UserModel{
		UUID:      *u.uuid,
		Username:  u.options["username"].(string),
		HPassword: u.options["password"].(string),
		Email:     u.options["email"].(string),
		Gender:    u.options["gender"].(model.UserGender),
		Region:    u.options["region"].(string),
		Other: model.UserOtherInfo{
			Introduction: u.options["other"].(map[string]any)["introduction"].(string),
			Icon:         u.options["other"].(map[string]any)["icon"].(string),
		},
	})
}

func (u *UserDAO) Update(ctx context.Context) error {
	//tx := u.db.Begin()
	//other, exists := u.options["other"]
	//if exists {
	//	otherMap := other.(map[string]any)
	//	for k, v := range otherMap {
	//		rowInfected, err := gorm.G[model.UserModel](tx).Update(ctx, "other", gorm.Expr("json_set(other, '$.?', ?)", k, v))
	//		if err != nil {
	//			tx.Rollback()
	//			return err
	//		}
	//		if rowInfected < 1 {
	//			tx.Rollback()
	//			return errors.New("nothing was updated")
	//		}
	//	}
	//}
	//// u.options["other"] = nil // 将other赋零值后调用updates会将整个json字段清空
	//delete(u.options, "other")
	//if err := tx.Model(&model.UserModel{}).Updates(u.options).Error; err != nil {
	//	tx.Rollback()
	//	return err
	//}
	//tx.Commit()
	//return nil
	// 多条sql语句性能低 合并为一条sql语句不用手动开启管理事务 gorm的占位符不能用于sql的路径 如 $.? 就是非法的
	if u.uuid == nil {
		return errors.New("uuid is nil")
	}
	if other, exists := u.options["other"]; exists {
		otherMap := other.(map[string]any)
		if len(otherMap) <= 0 {
			delete(u.options, "other") // 防止other为空值时更新把原本json数据清空
		} else {
			expr := "json_set(ifnull(other, '{}')"
			args := make([]any, len(otherMap)*2)
			for k, v := range otherMap {
				expr += ", ?, ?"
				args = append(args, "$."+k, v)
			}
			expr += ")"
			u.options["other"] = gorm.Expr(expr, args...)
		}
	}
	result := u.db.WithContext(ctx).Model(model.UserModel{}).Where("uuid = ?", *u.uuid).Updates(u.options)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return errors.New("no record was affected")
	}
	return nil
}

func (u *UserDAO) Delete(ctx context.Context) error {
	if u.uuid == nil {
		return errors.New("uuid is nil")
	}
	rowAffected, err := gorm.G[model.UserModel](u.db).Where("uuid = ?", *u.uuid).Delete(ctx)
	if err != nil {
		return err
	}
	if rowAffected < 1 {
		return errors.New("no record was affected")
	}
	return nil
}

func (u *UserDAO) Get(ctx context.Context) (*model.UserModel, error) {
	if u.uuid == nil {
		return nil, errors.New("uuid is nil")
	}
	first, err := gorm.G[model.UserModel](u.db).Where("uuid = ?", *u.uuid).First(ctx)
	if err != nil {
		return nil, err
	}
	return &first, nil
}

func (u *UserDAO) AddFollower(ctx context.Context, followerUUID string) error {
	tx := u.db.Begin()
	if u.uuid == nil {
		tx.Rollback()
		return errors.New("uuid is nil")
	}
	err := gorm.G[model.UserFollowers](tx).Create(ctx, &model.UserFollowers{
		FollowingID: *u.uuid,
		FollowerID:  followerUUID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	rowAffected, err := gorm.G[model.UserModel](tx).Where("uuid = ?", *u.uuid).Update(ctx, "follower_count", gorm.Expr("follower_count + 1"))
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowAffected < 1 {
		tx.Rollback()
		return errors.New("no record was affected")
	}
	rowAffected, err = gorm.G[model.UserModel](tx).Where("uuid = ?", followerUUID).Update(ctx, "following_count", gorm.Expr("following_count + 1"))
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowAffected < 1 {
		tx.Rollback()
		return errors.New("no record was affected")
	}
	tx.Commit()
	return nil
}

func (u *UserDAO) RemoveFollower(ctx context.Context, followerUUID string) error {
	tx := u.db.Begin()
	if u.uuid == nil {
		tx.Rollback()
		return errors.New("uuid is nil")
	}
	rowAffected, err := gorm.G[model.UserFollowers](tx).Where("follower_id = ? and following_id = ?", followerUUID, *u.uuid).Delete(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowAffected < 1 {
		tx.Rollback()
		return errors.New("no record was affected")
	}
	rowAffected, err = gorm.G[model.UserModel](tx).Where("uuid = ?", *u.uuid).Update(ctx, "follower_count", gorm.Expr("follower_count - 1"))
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowAffected < 1 {
		tx.Rollback()
		return errors.New("no record was affected")
	}
	rowAffected, err = gorm.G[model.UserModel](tx).Where("uuid = ?", followerUUID).Update(ctx, "following_count", gorm.Expr("following_count - 1"))
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowAffected < 1 {
		tx.Rollback()
		return errors.New("no record was affected")
	}
	tx.Commit()
	return nil
}

func (u *UserDAO) GetViaUsername(ctx context.Context, username string) (*model.UserModel, error) {
	user, err := gorm.G[model.UserModel](u.db).Where("username = ?", username).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDAO) SearchViaUsername(ctx context.Context, username string) ([]*model.UserModel, error) {
	// TODO
	return nil, errors.New("TODO")
}
