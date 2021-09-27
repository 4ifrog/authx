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
			setErrorStatus(ctx, ErrInvalidCredentials, http.StatusUnprocessableEntity)
			return
		}

		// Authenticate
		user, err := auth.Authenticate(ctx, ah.ds, login.Username, login.Password)
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidCredentials {
			setErrorStatus(ctx, ErrUserNotFound, http.StatusUnauthorized)
			return
		} else if err != nil {
			setErrorStatus(ctx, err, http.StatusInternalServerError)
			return
		}

		// Generate and save oauth2 object, which includes.
		aTTL := time.Duration(ah.cfg.AccessTTL) * time.Second
		rTTL := time.Duration(ah.cfg.RefreshTTL) * time.Second
		otoken, err := auth.CreateOAuthToken(ctx, ah.ds, user.ID, ah.cfg.AccessSecret, aTTL, rTTL)
		if err != nil {
			setErrorStatus(ctx, err, http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, otoken)
	}
}

func (ah *AuthHandlers) SignOut() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		atid := getAccessTokenIDFromContext(ctx)
		if atid == "" {
			ctx.JSON(http.StatusOK, "no access token")
			return
		}

		// Delete the token in the data store.
		claims, err := auth.UnsafeParseJWT(atid)
		if err != nil {
			setErrorStatus(ctx, err, http.StatusInternalServerError)
			return
		}
		if err := ah.ds.DeleteAccessToken(ctx, claims.ID); err != nil {
			setErrorStatus(ctx, err, http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, "logout")
	}
}

func (ah *AuthHandlers) UserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get(keyUserID)
		if !ok {
			setErrorStatus(ctx, ErrUserNotFound, http.StatusUnauthorized)
			return
		}

		uid, ok := val.(string)
		if !ok {
			setErrorStatus(ctx, ErrUserIDCast, http.StatusInternalServerError)
			return
		}

		user, err := ah.ds.GetUser(ctx, uid)
		if err == auth.ErrUserNotFound {
			setErrorStatus(ctx, ErrUserNotFound, http.StatusUnauthorized)
			return
		} else if err != nil {
			setErrorStatus(ctx, err, http.StatusInternalServerError)
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
			setErrorStatus(ctx, err, http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
		ctx.Writer.Header().Set("Content-Type", "image/svg+xml")
		ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(iconData)))
		_, err = ctx.Writer.Write(iconData)
		if err != nil {
			setErrorStatus(ctx, err, http.StatusInternalServerError)
			return
		}
	}
}
