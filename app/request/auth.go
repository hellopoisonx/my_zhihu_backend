package request

import "my_zhihu_backend/app/model"

type AuthLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthRenewAccessTokenRequest struct {
	Id           model.UserId `json:"user_id" binding:"required"`
	RefreshToken string       `json:"refresh_token" binding:"required"`
}
