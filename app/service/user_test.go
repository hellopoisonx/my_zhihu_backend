package service_test

import (
	"app/mock"
	"app/model"
	"app/mysql"
	"app/proxy"
	"app/request"
	"app/service"
	"app/utils"
	"context"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserDAO := mock.NewMockUserDAOProxy(ctrl)
	mockUserUtils := mock.NewMockUserUtilsProxy(ctrl)
	uuid := "test-uuid"
	req := request.UserCreateRequest{
		Username:     "test_user1",
		Password:     "test_password1",
		Email:        "test1@test.email",
		Gender:       utils.Ptr(model.UserGenderSecret),
		Region:       "test_region1",
		Introduction: "nothing here",
		Icon:         "http://test.icon",
	}
	hPassword := "123456"
	gomock.InOrder(
		mockUserUtils.EXPECT().EncryptPassword(req.Password).Return([]byte(hPassword), nil),
		mockUserUtils.EXPECT().GenerateUUID().Return(uuid),
		mockUserDAO.EXPECT().UUID(uuid).Return(mockUserDAO),
		mockUserDAO.EXPECT().Username(req.Username).Return(mockUserDAO),
		mockUserDAO.EXPECT().Password(hPassword).Return(mockUserDAO),
		mockUserDAO.EXPECT().Gender(req.Gender).Return(mockUserDAO),
		mockUserDAO.EXPECT().Region(req.Region).Return(mockUserDAO),
		mockUserDAO.EXPECT().Introduction(req.Introduction).Return(mockUserDAO),
		mockUserDAO.EXPECT().Icon(req.Icon).Return(mockUserDAO),
		// 最终执行动作
		mockUserDAO.EXPECT().Create(gomock.Any()).Return(nil),
	)
	var f = func(_ *mysql.Mysql) proxy.UserDAOProxy {
		return mockUserDAO
	}
	us := service.NewUserService(nil, f, mockUserUtils)
	err := us.CreateUser(context.TODO(), &req)

}
