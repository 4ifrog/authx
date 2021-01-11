package api

import (
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/storage"
)

func GetRoutesFunc() server.RegisterRoutesFunc {
	return func(parentGrp *gin.RouterGroup, cfg *config.Config, store storage.Storage) {
		parentGrp.POST("/signin", SignInHandler(cfg, store))
	}
}
