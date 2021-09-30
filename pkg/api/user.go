package api

import (
	"github.com/cybersamx/authx/pkg/models"
)

type SignInUser struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func user2UserInfo(user *models.User) *UserInfo {
	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
	}
}
