package request

import "my_zhihu_backend/app/model"

type CreateNewUserRequest struct {
	Username string               `json:"username" binding:"required"`
	Password string               `json:"password" binding:"required"`
	Email    string               `json:"email" binding:"required,email"`
	Region   string               `json:"region"`
	Gender   model.UserGender     `json:"gender" binding:"oneof=0 1 2"`
	Other    UserOtherInfoRequest `json:"other,omitempty"`
	Settings UserSettings         `json:"settings,omitempty"`
}

type UserOtherInfoRequest struct {
	Introduction *string `json:"introduction" binding:"omitempty"`
	Icon         *string `json:"icon" binding:"omitempty"`
}

type UserSettings struct {
	HidePrivacy *bool `json:"hide_privacy" binding:"omitempty"`
}

// UpdateUserRequest 用于更新用户信息
type UpdateUserRequest struct {
	Username string                `json:"username,omitempty"`
	Password string                `json:"password,omitempty"`
	Region   string                `json:"region,omitempty"`
	Gender   *model.UserGender     `json:"gender,omitempty" binding:"omitempty,oneof=0 1 2"` // 使用指针区分空值和零值
	Other    *UserOtherInfoRequest `json:"other,omitempty"`
	Settings *UserSettings         `json:"settings,omitempty"`
}

// SearchUserRequest 用于搜索用户
type SearchUserRequest struct {
	Username string `form:"username" binding:"required"`
}
