package api

import (
	"github.com/cybersamx/authx/pkg/models"
)

type SignInRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserInfoResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func user2UserInfoResponse(user *models.User) *UserInfoResponse {
	return &UserInfoResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}
