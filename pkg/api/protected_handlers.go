package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
)

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func userToUserInfo(user *models.User) *UserInfo {
	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
	}
}

func UserInfoHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		obj, ok := ctx.Get("UserID")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, ok := obj.(string)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err := ds.GetUser(ctx, userID)
		if err == auth.ErrUserNotFound {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		} else if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, userToUserInfo(user))
	}
}
