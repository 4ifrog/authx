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
		htmlGrp := router.Group("/")
		htmlGrp.GET("/", HomeHTMLHandler(cfg, ds))
		htmlGrp.GET("/signin", SignInHTMLHandler(cfg, ds))
		htmlGrp.POST("/signin", SignInFormHTMLHandler(cfg, ds))
		htmlGrp.GET("/profile", ProfileHTMLHandler(cfg, ds))
		htmlGrp.POST("/signout", SignoutFormHTMLHandler(cfg, ds))

		apiGrp := router.Group("/v1")
		apiGrp.POST("/signin", SignInHandler(cfg, ds))
		apiGrp.GET("/signout", SignOutHandler(cfg, ds))

		// Place this last.
		router.Use(static.ServeRoot("/", "./public"))
	}
}
