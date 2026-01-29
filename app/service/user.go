package service

import (
	"app/app_error"
	"app/mysql"
	"app/proxy"
	"app/request"
	"app/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrUserAlreadyExists = app_error.NewInputError(
	"user already exists", 201, nil,
)

type UserDAOFactory func(*mysql.Mysql) proxy.UserDAOProxy

type UserService struct {
	mysql   *mysql.Mysql
	factory UserDAOFactory
	utils   proxy.UserUtilsProxy
}

func NewUserService(mysql *mysql.Mysql, factory UserDAOFactory, utils proxy.UserUtilsProxy) *UserService {
	return &UserService{
		mysql:   mysql,
		factory: factory,
		utils:   utils,
	}
}

func (u *UserService) CreateUser(ctx context.Context, request *request.UserCreateRequest) error {
	hPassword, err := u.utils.EncryptPassword(request.Password)
	if err != nil {
		return app_error.NewInternalError(utils.ErrCodeEncryptionFailed, err)
	}
	err = u.factory(u.mysql).
		UUID(u.utils.GenerateUUID()).
		Username(request.Username).
		Password(string(hPassword)).
		Email(request.Email).
		Gender(*request.Gender).
		Region(request.Region).
		Introduction(request.Introduction).
		Icon(request.Icon).
		Create(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUserAlreadyExists
		}
		return app_error.NewInternalError(mysql.ErrCodeMysql, err)
	}
	return nil
}
