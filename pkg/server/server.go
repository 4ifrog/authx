package server

import (
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/store"
)

type RegisterRoutesFunc func(parentGrp *gin.RouterGroup, cfg *config.Config, store store.DataStore)

type Server struct {
	Router  *gin.Engine
	rootGrp *gin.RouterGroup
	cfg     *config.Config
}

func (s *Server) BindAPIRoutes(fn RegisterRoutesFunc, ds store.DataStore) {
	fn(s.rootGrp, s.cfg, ds)
}

func New(cfg *config.Config) *Server {
	s := &Server{
		Router: gin.New(),
		cfg:    cfg,
	}

	s.rootGrp = s.Router.Group("/v1")

	return s
}
