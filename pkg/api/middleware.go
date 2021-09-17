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

// getUserIDFromContext gets the user ID from the context.
func getUserIDFromContext(ctx *gin.Context) string {
	userID, ok := ctx.Get(keyUserID)
	if !ok {
		return ""
	}

	return userID.(string)
}

type Middleware struct {
	cfg *config.Config
	ds  store.DataStore
}

func NewMiddleware(cfg *config.Config, ds store.DataStore) *Middleware {
	mw := new(Middleware)

	mw.cfg = cfg
	mw.ds = ds

	return mw
}

// AccessTokenFromBearerAuth identifies the user of a request by extracting the user id from the JWT of
// the request before setting the user id in the context.
func (m *Middleware) AccessTokenFromBearerAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// See the user id has already been set.
		uid := getUserIDFromContext(ctx)
		if uid != "" {
			// Context already has the user id.
			return
		}

		// Extract the bearer token from the header and parse it as JWT.
		bearerToken, err := parseBearerFromHeader(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			return
		}
		at, err := auth.ParseJWT(bearerToken, m.cfg.AccessSecret)
		if err != nil || at == nil {
			return
		}

		// TODO: Check if user exists in the data store.
		ctx.Set(keyUserID, at.UserID)
		ctx.Set(keyAccessTokenID, at.ID)

		ctx.Next()
	}
}

func (m *Middleware) UserFromCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// See the user id has already been set.
		uid := getUserIDFromContext(ctx)
		if uid != "" {
			// Context already has the user id.
			return
		}

		// Extract user id and access token from the cookie.
		ss := NewSessionStore(m.cfg.SessionSecret)
		us, err := ss.GetSession(ctx.Request)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if us == nil {
			return
		}

		// Check that user exists in the data store.
		user, err := m.ds.GetUser(ctx, us.UserID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Set(keyUserID, user.ID)
		ctx.Set(keyAccessTokenID, us.OAuth2Token.AccessToken)

		ctx.Next()
	}
}

func (m *Middleware) AccessTokenFromCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session
		ss := NewSessionStore(m.cfg.SessionSecret)
		session, err := ss.GetSession(ctx.Request)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		at := session.OAuth2Token.AccessToken
		claims, err := auth.ParseJWT(at, m.cfg.AccessSecret)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Set("UserID", claims.UserID)
		ctx.Next()
	}
}
