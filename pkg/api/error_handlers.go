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
	errorsPathFormat = "/errors/%d"
)

var (
	statusInternalServerContent = gin.H{
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

func (eh *ErrorHandlers) ErrorResponse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Run other handlers in the chain first.
		ctx.Next()

		// Ignore 2xx and 3xx status code.
		code := ctx.Writer.Status()
		if code >= http.StatusOK && code < http.StatusBadRequest {
			return
		}

		if isAcceptingHTML(ctx.Request) {
			ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf(errorsPathFormat, code))
			ctx.Abort()
			return
		}

		ctx.AbortWithStatus(code)
	}
}

func (eh *ErrorHandlers) ErrorHTML() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code, err := strconv.Atoi(ctx.Param("code"))
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "error", statusInternalServerContent)
			ctx.Abort()
			return
		}

		content := gin.H{
			"Code":       code,
			"StatusText": http.StatusText(code),
		}

		ctx.HTML(code, "error", content)
	}
}
