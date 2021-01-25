package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/authn"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/models"
	"github.com/cybersamx/authx/pkg/storage"
)

const (
	bearerComponentCount = 2
)

type tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func getBearerFromRequest(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	components := strings.Split(authHeader, " ")
	if len(components) == bearerComponentCount {
		return components[1]
	}

	return ""
}

func unsignedJWT(accessSecret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported jwt signing method: %s", token.Header["alg"])
		}

		return []byte(accessSecret), nil
	}
}

func ValidateToken(req *http.Request, accessSecret string) (*jwt.Token, error) {
	bearer := getBearerFromRequest(req)
	token, err := jwt.Parse(bearer, unsignedJWT(accessSecret))

	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetToken(req *http.Request, accessSecret string) error {
	token, err := ValidateToken(req, accessSecret)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.Claims)
	if !ok && !token.Valid {
		return err
	}

	fmt.Println(claims)

	return nil
}

func SignInHandler(cfg *config.Config, store storage.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var login models.User
		if err := ctx.ShouldBindJSON(&login); err != nil {
			ctx.JSONP(http.StatusUnprocessableEntity, "Invalid request payload")
			return
		}

		user, err := authn.Authenticate(ctx, store, login.Username, login.Password)
		if err == authn.ErrUserNotFound || err == authn.ErrInvalidCredentials {
			ctx.JSONP(http.StatusUnauthorized, "Invalid authentication credentials")
			return
		} else if err != nil {
			ctx.JSONP(http.StatusInternalServerError, err.Error())
			return
		}

		at, rt, err := createOAuthToken(user.ID, cfg)
		if err != nil {
			ctx.JSONP(http.StatusInternalServerError, err.Error())
			return
		}

		if err := saveOAuthToken(ctx, store, at, rt); err != nil {
			ctx.JSONP(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSONP(http.StatusOK, tokens{
			AccessToken:  at.Value,
			RefreshToken: rt.Value,
		})
	}
}

func SignOutHandler(cfg *config.Config, store storage.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSONP(http.StatusNotImplemented, "Not implemented")
	}
}
