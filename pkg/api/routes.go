package api

import (
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/store"
)

func GetRoutesFunc() server.RegisterRoutesFunc {
	return func(router *gin.Engine, cfg *config.Config, ds store.DataStore) {
		// Auth API.
		apiGrp := router.Group("/v1")
		apiGrp.POST("/signin", SignInHandler(cfg, ds))
		apiGrp.GET("/signout", SignOutHandler(cfg, ds))

		// React SPA.
		router.Use(StaticHandler(cfg))
	}
}
