package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/avatar"
)

// AvatarHandler returns identicon avatar icon.
func AvatarHandler() gin.HandlerFunc {
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
