package service

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/dao"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/request"
	"my_zhihu_backend/app/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "this is a secret key" // TODO: 更换密钥

type UserJWTClaims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	aDAO *dao.AuthDAO
	uDAO *dao.UserDAO
	cfg  config.ReadConfigFunc
	util *util.Util
}

func NewAuthService(aDAO *dao.AuthDAO, uDAO *dao.UserDAO, cfg config.ReadConfigFunc, util *util.Util) *AuthService {
	return &AuthService{aDAO, uDAO, cfg, util}
}

func (s *AuthService) newAccessToken(id model.UserId) (string, time.Time, app_error.AppError) {
	expAt := time.Now().Add(s.cfg().Service.AccessTokenExp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserJWTClaims{
		Id: int64(id),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
	})
	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Now(), app_error.NewInternalError(app_error.ErrCodeUserToken, err)
	}
	return signedString, expAt, nil
}

func (s *AuthService) ValidateAccessToken(token string) (model.UserId, app_error.AppError) {
	t, err := jwt.ParseWithClaims(token, &UserJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, app_error.ErrUserInvalidToken.WithError(err)
	}
	if !t.Valid {
		return 0, app_error.ErrUserInvalidToken
	}
	claims := t.Claims.(*UserJWTClaims)
	return model.UserId(claims.Id), nil
}

func (s *AuthService) Login(ctx context.Context, req *request.AuthLoginRequest) (accessToken, refreshToken string, accessExpireAt, refreshExpireAt time.Time, err app_error.AppError) {
	var user *model.User
	user, err = s.uDAO.GetByEmail(ctx, req.Email)
	if err != nil {
		return
	}
	if !s.util.ValidatePassword(user.HPassword, req.Password) {
		err = app_error.ErrUserWrongPassword
		return
	}
	accessToken, accessExpireAt, err = s.newAccessToken(user.Id)
	if err != nil {
		return
	}
	refreshToken = s.util.GenerateUUID()
	refreshExpireAt = time.Now().Add(s.cfg().Service.RefreshTokenExp)
	if err = s.aDAO.SaveRefreshToken(ctx, user.Id, refreshToken, s.cfg().Service.RefreshTokenExp); err != nil {
		return
	}
	return
}

func (s *AuthService) Logout(ctx context.Context, id model.UserId) app_error.AppError {
	return s.aDAO.DeleteRefreshToken(ctx, id)
}

func (s *AuthService) RenewAccessToken(ctx context.Context, req *request.AuthRenewAccessTokenRequest) (string, time.Time, app_error.AppError) {
	token, err := s.aDAO.GetRefreshToken(ctx, req.Id)
	if err != nil {
		return "", time.Now(), err
	}
	if token != req.RefreshToken {
		return "", time.Now(), app_error.ErrUserInvalidToken
	}
	return s.newAccessToken(req.Id)
}
