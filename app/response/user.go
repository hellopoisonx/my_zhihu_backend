package response

import "my_zhihu_backend/app/model"

type UserResponse struct {
	Id       model.UserId          `json:"id"`
	Username string                `json:"username"`
	Email    string                `json:"email"`
	Gender   model.UserGender      `json:"gender"`
	Region   string                `json:"region"`
	Other    UserOtherInfoResponse `json:"other"`
}

type UserOtherInfoResponse struct {
	Introduction string `json:"introduction"`
	Icon         string `json:"icon"`
}
