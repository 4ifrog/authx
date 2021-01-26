package api

import (
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/store"
)

func GetRoutesFunc() server.RegisterRoutesFunc {
	return func(parentGrp *gin.RouterGroup, cfg *config.Config, ds store.DataStore) {
		parentGrp.POST("/signin", SignInHandler(cfg, ds))
		parentGrp.GET("/signout", SignOutHandler(cfg, ds))
	}
}
