package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cybersamx/authx/pkg/auth"
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/store"
)

const (
	keyUserID        = "UserID"
	keyAccessTokenID = "AccessTokenID"
)

var (
	ErrMissingBearer        = errors.New("missing bearer")
	ErrInvalidBearer        = errors.New("bearer has invalid content")
	ErrMissingSessionCookie = errors.New("missing session cookie")
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

// getAccessTokenIDFromContext gets the access token from the context.
func getAccessTokenIDFromContext(ctx *gin.Context) string {
	at, ok := ctx.Get(keyAccessTokenID)
	if !ok {
		return ""
	}

	return at.(string)
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

// TODO: Consider refactoring, too many return params.
type extractFunc func(ctx *gin.Context) (string, string, int, error)

func (m *Middleware) setContextUsing(ctx *gin.Context, extractor extractFunc) {
	// See the user id has already been set.
	uid := getUserIDFromContext(ctx)
	if uid != "" {
		// Context already has the user id.
		ctx.Next()
		return
	}

	// Extracts user id and access token id using the extractor.
	uid, atid, status, err := extractor(ctx)
	if err != nil {
		_ = ctx.AbortWithError(status, err)
		return
	}

	// Check if user is in the data store.
	_, err = m.ds.GetUser(ctx, uid)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	} else if err == store.ErrorNotFound {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Set the access token id and user in the context for other handlers to access.
	ctx.Set(keyUserID, uid)
	ctx.Set(keyAccessTokenID, atid)

	ctx.Next()
}

// SetContextFromBearerAuth extracts the user identity from the bearer token and if it's valid, saves
// it to the context for other handlers to use.
func (m *Middleware) SetContextFromBearerAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fun := func(ctx *gin.Context) (string, string, int, error) {
			// Extract the bearer token from the header and parse it as JWT.
			bearerToken, err := parseBearerFromHeader(ctx.Request.Header.Get("Authorization"))
			if err != nil {
				return "", "", http.StatusInternalServerError, err
			}
			at, err := auth.ParseJWT(bearerToken, m.cfg.AccessSecret)
			if err == auth.ErrExpiredJWT {
				return "", "", http.StatusUnauthorized, err
			} else if err != nil {
				return "", "", http.StatusInternalServerError, err
			}

			return at.UserID, at.ID, http.StatusOK, nil
		}

		m.setContextUsing(ctx, fun)
	}
}

// SetContextFromCookie extracts the user identity from the session cookie and if it's valid, saves
// it to the context for other handlers to use.
func (m *Middleware) SetContextFromCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fun := func(ctx *gin.Context) (string, string, int, error) {
			ss := NewCookieStore(m.cfg.SessionSecret)
			us, err := ss.GetSessionToken(ctx.Request)
			if err == ErrSessionTokenNotFound {
				return "", "", http.StatusUnauthorized, ErrMissingSessionCookie
			} else if err != nil {
				return "", "", http.StatusInternalServerError, err
			}

			if us == nil {
				return "", "", http.StatusUnauthorized, ErrMissingSessionCookie
			}

			return us.UserID, us.Token.AccessToken, http.StatusOK, nil
		}

		m.setContextUsing(ctx, fun)
	}
}
