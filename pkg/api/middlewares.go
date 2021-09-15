package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/store"
)

const (
	keyUserID        = "UserID"
	keyAccessTokenID = "AccessTokenID"
)

var (
	ErrMissingBearer = errors.New("missing bearer")
	ErrInvalidBearer = errors.New("bearer has invalid content")
)

// parseBearerFromHeader gets the bearer token from the header string.
func parseBearerFromHeader(header string) (string, error) {
	if header == "" {
		return "", ErrMissingBearer
	}

	chunks := strings.Split(header, " ")
	if len(chunks) == 2 && chunks[0] == "Bearer" {
		return chunks[1], nil
	}

	return "", ErrInvalidBearer
}

// GetUserIDFromContext gets the user ID from the context.
func GetUserIDFromContext(ctx *gin.Context) string {
	userID, ok := ctx.Get(keyUserID)
	if !ok {
		return ""
	}

	return userID.(string)
}

// GetAccessTokenIDFromContext gets the access token ID from the context.
func GetAccessTokenIDFromContext(ctx *gin.Context) string {
	atID, ok := ctx.Get(keyAccessTokenID)
	if !ok {
		return ""
	}

	return atID.(string)
}

// BearerAuthHandler identifies the user of a request by extracting the user id from the JWT of
// the request before setting the user id in the context.
func BearerAuthHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// See the user id has already been set.
		uid := GetUserIDFromContext(ctx)
		if uid != "" {
			// Context already has the user id.
			return
		}

		// Extract the bearer token from the header and parse it as JWT.
		bearerToken, err := parseBearerFromHeader(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			return
		}
		at, err := auth.ParseJWT(bearerToken, cfg.AccessSecret)
		if err != nil || at == nil {
			return
		}

		// TODO: Check if user exists in the data store.
		ctx.Set(keyUserID, at.UserID)
		ctx.Set(keyAccessTokenID, at.ID)

		ctx.Next()
	}
}

func CheckAuthorizationHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		at := GetAccessTokenIDFromContext(ctx)
		if at == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Next()
	}
}

func CookieHandler(cfg *config.Config, ds store.DataStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// See the user id has already been set.
		uid := GetUserIDFromContext(ctx)
		if uid != "" {
			// Context already has the user id.
			return
		}

		// Extract user id and access token from the cookie.
		ss := NewSessionStore(cfg.SessionSecret)
		us, err := ss.GetSession(ctx.Request)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if us == nil {
			return
		}

		// Check that user exists in the data store.
		user, err := ds.GetUser(ctx, us.UserID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Set(keyUserID, user.ID)
		ctx.Set(keyAccessTokenID, us.OAuth2Token.AccessToken)

		ctx.Next()
	}
}
