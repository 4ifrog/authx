package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/cybersamx/authx/pkg/avatar"
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

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AuthHandlers struct {
	cfg *config.Config
	ds  store.DataStore
}

func NewAuthHandlers(cfg *config.Config, ds store.DataStore) *AuthHandlers {
	handlers := new(AuthHandlers)

	handlers.cfg = cfg
	handlers.ds = ds

	return handlers
}

func (ah *AuthHandlers) userToUserInfo(user *models.User) *UserInfo {
	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
	}
}

func (ah *AuthHandlers) SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind inputs
		var login models.User
		if err := ctx.ShouldBindJSON(&login); err != nil {
			_ = ctx.AbortWithError(http.StatusUnprocessableEntity, ErrInvalidCredentials)
			return
		}

		// Authenticate
		user, err := auth.Authenticate(ctx, ah.ds, login.Username, login.Password)
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidCredentials {
			_ = ctx.AbortWithError(http.StatusUnauthorized, ErrInvalidRequest)
			return
		} else if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Generate oauth2 object and save.
		aTTL := time.Duration(ah.cfg.AccessTTL) * time.Second
		rTTL := time.Duration(ah.cfg.RefreshTTL) * time.Second
		otoken, err := auth.CreateOAuthToken(ctx, ah.ds, user.ID, ah.cfg.AccessSecret, aTTL, rTTL)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Save session to the cookie.
		session := UserSession{
			OAuth2Token: *otoken,
			UserID:      user.ID,
		}
		ss := NewSessionStore(ah.cfg.SessionSecret)
		if err := ss.SetSession(ctx.Writer, ctx.Request, &session); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, otoken)
	}
}

func (ah *AuthHandlers) SignOut() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session
		ss := NewSessionStore(ah.cfg.SessionSecret)
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
		if err := ah.ds.DeleteAccessToken(ctx, claims.ID); err != nil {
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

func (ah *AuthHandlers) UserInfo() gin.HandlerFunc {
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

		user, err := ah.ds.GetUser(ctx, userID)
		if err == auth.ErrUserNotFound {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		} else if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, ah.userToUserInfo(user))
	}
}

// AvatarHandler returns identicon avatar icon.
func (ah *AuthHandlers) Avatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		identity := ctx.Param("identity")
		iconData, err := avatar.GetIdenticon(identity)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Status(http.StatusOK)
		ctx.Writer.Header().Set("Content-Type", "image/svg+xml")
		ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(iconData)))
		_, err = ctx.Writer.Write(iconData)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
