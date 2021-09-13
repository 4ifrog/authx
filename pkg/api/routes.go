package api

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/store"
)

func GetRoutesFunc() server.RegisterRoutesFunc {
	return func(router *gin.Engine, cfg *config.Config, ds store.DataStore) {
		// Web Pages.
		webGrp := router.Group("/")
		webGrp.GET("/", WebSignInHandler(cfg, ds))
		webGrp.POST("/", WebSignInHandler(cfg, ds))

		// Auth Public API.
		apiGrp := router.Group("/v1")
		apiGrp.POST("/signin", SignInHandler(cfg, ds))
		apiGrp.GET("/signout", SignOutHandler(cfg, ds))
		apiGrp.GET("/avatar/:identity", AvatarHandler())

		// Protected Auth API.
		protectedGrp := router.Group("/v1")
		protectedGrp.Use(AccessTokenHandler(cfg))
		protectedGrp.GET("/userinfo", UserInfoHandler(cfg, ds))

		// Fallback to static content.
		router.Use(static.Serve("/", static.LocalFile(cfg.StaticWebDir, false)))
	}
}
