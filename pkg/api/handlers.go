package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/auth/oauth2"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
)

func SignInHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var login models.User
		if err := ctx.ShouldBindJSON(&login); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, "Invalid request payload")
			return
		}

		user, err := auth.Authenticate(ctx, ds, login.Username, login.Password)
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidCredentials {
			ctx.JSON(http.StatusUnauthorized, "Invalid authentication credentials")
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		accessTTL := time.Duration(cfg.AccessTTL) * time.Second
		refreshTTL := time.Duration(cfg.RefreshTTL) * time.Second
		at, rt, err := oauth2.NewOAuthToken(user.ID, cfg.AccessSecret, accessTTL, refreshTTL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if err := oauth2.SaveOAuthToken(ctx, ds, at, rt); err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, OAuth2Tokens{
			AccessToken:  NewAccessTokenModel(at),
			RefreshToken: NewRefreshTokenModel(rt),
		})
	}
}

func SignOutHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotImplemented, "Not implemented")
	}
}
