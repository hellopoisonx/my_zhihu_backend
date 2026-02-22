package controller

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/request"
	"my_zhihu_backend/app/response"
	"my_zhihu_backend/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
	cfg     config.ReadConfigFunc
}

func NewUserController(service *service.UserService, cfg config.ReadConfigFunc) *UserController {
	return &UserController{service: service, cfg: cfg}
}

func (ctrl *UserController) CreateNewUser(c *gin.Context) {
	do(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, req *request.CreateNewUserRequest) (*response.Response, app_error.AppError) {
		if user, err := ctrl.service.CreateNewUser(ctx, req); err != nil {
			return nil, err
		} else {
			return &response.Response{
				Ok:            true,
				InternalError: false,
				Code:          0,
				Message:       "user created",
				Body: response.UserResponse{
					Id:       user.Id,
					Username: user.Username,
					Email:    user.Email,
					Gender:   *user.Gender,
					Region:   user.Region,
					Other: response.UserOtherInfoResponse{
						Introduction: user.Other.Introduction,
						Icon:         user.Other.Icon,
					},
				},
			}, nil
		}
	})
}

// DeleteUser 删除用户
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ctrl.cfg().Service.Timeout)
	defer cancel()
	id := getCurrentUserID(c)

	if err := ctrl.service.DeleteUser(timeout, id); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, &response.Response{
		Ok:            true,
		InternalError: false,
		Code:          0,
		Message:       "user deleted",
	})
}

// GetUser 获取用户信息
func (ctrl *UserController) GetUser(c *gin.Context) {
	doQuery(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, req *struct{}) (*response.Response, app_error.AppError) {
		id, err := getUserIDFromParam(c)
		if err != nil {
			return nil, app_error.NewInputError("invalid user id", app_error.ErrCodeInvalidParameters, err)
		}

		user, err := ctrl.service.GetUser(ctx, id)
		if err != nil {
			return nil, err.(app_error.AppError)
		}

		resp := response.UserResponse{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
			Gender:   *user.Gender,
			Region:   user.Region,
			Other: response.UserOtherInfoResponse{
				Introduction: user.Other.Introduction,
				Icon:         user.Other.Icon,
			},
		}

		if user.Settings.HidePrivacy && user.Id != model.UserId(getCurrentUserID(c)) {
			resp.Email = "nothing here"
			resp.Gender = model.UserGenderSecret
			resp.Region = "nothing here"
		}

		return &response.Response{
			Ok:            true,
			InternalError: false,
			Code:          0,
			Message:       "user retrieved",
			Body:          &resp,
		}, nil
	})
}

// UpdateUser 更新用户信息
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	do(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, req *request.UpdateUserRequest) (*response.Response, app_error.AppError) {
		currentUserID := getCurrentUserID(c)

		if user, err := ctrl.service.UpdateUser(ctx, currentUserID, req); err != nil {
			return nil, err
		} else {
			return &response.Response{
				Ok:            true,
				InternalError: false,
				Code:          0,
				Message:       "user updated",
				Body: response.UserResponse{
					Id:       user.Id,
					Username: user.Username,
					Email:    user.Email,
					Gender:   *user.Gender,
					Region:   user.Region,
					Other: response.UserOtherInfoResponse{
						Introduction: user.Other.Introduction,
						Icon:         user.Other.Icon,
					},
				},
			}, nil
		}
	})
}

// SearchUserByUsername 根据用户名搜索用户
func (ctrl *UserController) SearchUserByUsername(c *gin.Context) {
	doQuery(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, req *request.SearchUserRequest) (*response.Response, app_error.AppError) {
		userIDs, err := ctrl.service.SearchUserByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		return &response.Response{
			Ok:            true,
			InternalError: false,
			Code:          0,
			Message:       "users found",
			Body:          userIDs,
		}, nil
	})
}

// AddFollowing 关注用户
func (ctrl *UserController) AddFollowing(c *gin.Context) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ctrl.cfg().Service.Timeout)
	defer cancel()
	currentUserID := getCurrentUserID(c)
	followingID, err := getUserIDFromParam(c)
	if err != nil {
		_ = c.Error(app_error.NewInputError("invalid user id", app_error.ErrCodeInvalidParameters, err))
		return
	}

	if err := ctrl.service.FollowUser(timeout, currentUserID, followingID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, &response.Response{
		Ok:            true,
		InternalError: false,
		Code:          0,
		Message:       "user followed",
	})
}

// RemoveFollowing 取消关注用户
func (ctrl *UserController) RemoveFollowing(c *gin.Context) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ctrl.cfg().Service.Timeout)
	defer cancel()
	currentUserID := getCurrentUserID(c)
	followingID, err := getUserIDFromParam(c)
	if err != nil {
		_ = c.Error(app_error.NewInputError("invalid user id", app_error.ErrCodeInvalidParameters, err))
		return
	}

	if err := ctrl.service.UnfollowUser(timeout, currentUserID, followingID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, &response.Response{
		Ok:            true,
		InternalError: false,
		Code:          0,
		Message:       "user unfollowed",
	})
}

// GetFollowers 获取粉丝列表
func (ctrl *UserController) GetFollowers(c *gin.Context) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ctrl.cfg().Service.Timeout)
	defer cancel()
	id, err := getUserIDFromParam(c)
	if err != nil {
		_ = c.Error(app_error.NewInputError("invalid parameters", app_error.ErrCodeInvalidParameters, err))
		return
	}
	currentUserID := getCurrentUserID(c)

	if currentUserID != id {
		user, err := ctrl.service.GetUser(timeout, id)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if user.Settings.HidePrivacy {
			_ = c.Error(app_error.ErrUserPermissionDenied)
			return
		}
	}
	followers, err := ctrl.service.GetFollowers(timeout, model.UserId(id))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Ok:            true,
		InternalError: false,
		Code:          0,
		Message:       "user followers",
		Body:          followers,
	})
}

// GetFollowings 获取关注列表
func (ctrl *UserController) GetFollowings(c *gin.Context) {
	timeout, cancel := context.WithTimeout(c.Request.Context(), ctrl.cfg().Service.Timeout)
	defer cancel()
	id, err := getUserIDFromParam(c)
	if err != nil {
		_ = c.Error(app_error.NewInputError("invalid parameters", app_error.ErrCodeInvalidParameters, err))
		return
	}
	currentUserID := getCurrentUserID(c)

	if currentUserID != id {
		user, err := ctrl.service.GetUser(timeout, id)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if user.Settings.HidePrivacy {
			_ = c.Error(app_error.ErrUserPermissionDenied)
			return
		}
	}
	followings, err := ctrl.service.GetFollowings(timeout, model.UserId(id))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &response.Response{
		Ok:            true,
		InternalError: false,
		Code:          0,
		Message:       "user followings",
		Body:          followings,
	})
}

// 辅助函数：从路径参数中获取用户ID
func getUserIDFromParam(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

// 辅助函数：从上下文获取当前用户ID
func getCurrentUserID(c *gin.Context) int64 {
	userID := c.MustGet("id").(model.UserId)

	return int64(userID)
}
