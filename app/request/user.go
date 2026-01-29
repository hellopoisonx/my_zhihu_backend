package request

import "app/model"

type UserCreateRequest struct {
	Username string            `json:"username" binding:"required,min=1,max=50"`
	Password string            `json:"password" binding:"required"`
	Email    string            `json:"email" binding:"required,email"`
	Gender   *model.UserGender `json:"gender" binding:"required,oneof=0 1 2"` // 不使用指针 gin绑定0到gender的时候会报错 required标签不允许零值

	Region       string `json:"region" binding:"required;max=50"`
	Introduction string `json:"introduction" binding:"max=50"`
	Icon         string `json:"icon" binding:"url"`
}
