//go:generate mockgen -source=user.go -destination=../mock/mock_user.go -package=mock
package proxy

import (
	"app/model"
	"app/request"
	"context"
)

type UserDAOProxy interface {
	UUID(string) UserDAOProxy
	Username(string) UserDAOProxy
	Password(string) UserDAOProxy
	Email(string) UserDAOProxy
	Gender(model.UserGender) UserDAOProxy
	Region(string) UserDAOProxy
	Introduction(string) UserDAOProxy
	Icon(string) UserDAOProxy
	Update(ctx context.Context) error
	Create(ctx context.Context) error
	Get(ctx context.Context) (*model.UserModel, error)
	Delete(ctx context.Context) error
	AddFollower(ctx context.Context, followerUUID string) error
	RemoveFollower(ctx context.Context, followerUUID string) error
	SearchViaUsername(ctx context.Context, username string) ([]*model.UserModel, error)
	GetViaUsername(ctx context.Context, username string) (*model.UserModel, error)
}

type UserServiceProxy interface {
	CreateUser(context.Context, *request.UserCreateRequest) error
}

type UserUtilsProxy interface {
	GenerateUUID() string
	EncryptPassword(string) ([]byte, error)
	ValidatePassword(hPassword, password string) bool
}
