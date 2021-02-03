package api

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
)

func StaticHandler(cfg *config.Config) gin.HandlerFunc {
	fs := static.LocalFile(cfg.StaticWebDir, true)
	fileserver := http.StripPrefix("/", http.FileServer(fs))

	return func(ctx *gin.Context) {
		if fs.Exists("/", ctx.Request.URL.Path) {
			// For serving html, js, css, and other assets that exist on the server.
			fileserver.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
		} else {
			// The React app uses client-side routing (spa), the route path doesn't
			// on the server so we force the root (of the spa) to be served.
			defer func(old string) {
				ctx.Request.URL.Path = old
			}(ctx.Request.URL.Path)

			ctx.Request.URL.Path = "/"
			fileserver.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Next()
		}
	}
}
