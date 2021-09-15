package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/store"
)

var (
	ErrInvalidCredentials = errors.New("invalid authentication credentials")
	ErrInvalidRequest     = errors.New("invalid request payload")
)

func SignInHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind inputs
		var login models.User
		if err := ctx.ShouldBindJSON(&login); err != nil {
			_ = ctx.AbortWithError(http.StatusUnprocessableEntity, ErrInvalidCredentials)
			return
		}

		// Authenticate
		user, err := auth.Authenticate(ctx, ds, login.Username, login.Password)
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidCredentials {
			_ = ctx.AbortWithError(http.StatusUnauthorized, ErrInvalidRequest)
			return
		} else if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Generate oauth2 object and save.
		aTTL := time.Duration(cfg.AccessTTL) * time.Second
		rTTL := time.Duration(cfg.RefreshTTL) * time.Second
		otoken, err := auth.CreateOAuthToken(ctx, ds, user.ID, cfg.AccessSecret, aTTL, rTTL)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Save session to the cookie.
		session := UserSession{
			OAuth2Token: *otoken,
			UserID:      user.ID,
		}
		ss := NewSessionStore(cfg.SessionSecret)
		if err := ss.SetSession(ctx.Writer, ctx.Request, &session); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, otoken)
	}
}

func SignOutHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session
		ss := NewSessionStore(cfg.SessionSecret)
		session, err := ss.GetSession(ctx.Request)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if session == nil {
			ctx.JSON(http.StatusOK, "logout")
			return
		}

		// Delete the token in the data store.
		claims, err := auth.UnsafeParseJWT(session.OAuth2Token.AccessToken)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if err := ds.DeleteAccessToken(ctx, claims.ID); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Clear session
		if err := ss.ClearSession(ctx.Writer, ctx.Request); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, "logout")
	}
}

func AccessTokenHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session
		ss := NewSessionStore(cfg.SessionSecret)
		session, err := ss.GetSession(ctx.Request)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		at := session.OAuth2Token.AccessToken
		claims, err := auth.ParseJWT(at, cfg.AccessSecret)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Set("UserID", claims.UserID)
		ctx.Next()
	}
}
