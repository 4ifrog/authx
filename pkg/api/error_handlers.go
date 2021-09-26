package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cybersamx/authx/pkg/config"

	"github.com/gin-gonic/gin"
)

const (
	errorsPath       = "/errors/:code"
	errorsPathFormat = "/errors/%v"
)

var (
	statusInternalErrorContent = gin.H{
		"Code":       http.StatusInternalServerError,
		"StatusText": http.StatusText(http.StatusInternalServerError),
	}
)

type ErrorHandlers struct {
	cfg *config.Config
}

func NewErrorHandlers(cfg *config.Config) *ErrorHandlers {
	handlers := new(ErrorHandlers)

	handlers.cfg = cfg

	return handlers
}

func (eh *ErrorHandlers) AbortWithError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, found := ctx.Get(keyStatusCode)
		if found {
			code, ok := val.(int)
			if !ok {
				_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't cast status code"))
				return
			}

			if code >= http.StatusOK && code < http.StatusMultipleChoices {
				// Ignore 2xx status code.
				ctx.Next()
				return
			}

			ctx.AbortWithStatus(code)
		}

		ctx.Next()
	}
}

func (eh *ErrorHandlers) RedirectToErrorPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, found := ctx.Get(keyStatusCode)
		if found {
			code, ok := val.(int)
			if !ok {
				ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf(errorsPathFormat, http.StatusInternalServerError))
				ctx.Abort()
				return
			}

			if code >= http.StatusOK && code < http.StatusMultipleChoices {
				// Ignore 2xx status code.
				ctx.Next()
				return
			}

			ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf(errorsPathFormat, val))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (eh *ErrorHandlers) Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code, err := strconv.Atoi(ctx.Param("code"))
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "error", statusInternalErrorContent)
			ctx.Abort()
			return
		}

		content := gin.H{
			"Code":       code,
			"StatusText": http.StatusText(code),
		}

		ctx.HTML(http.StatusUnauthorized, "error", content)
	}
}
