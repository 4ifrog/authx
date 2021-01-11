package server

import (
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/storage"
)

type RegisterRoutesFunc func(parentGrp *gin.RouterGroup, cfg *config.Config, store storage.Storage)

type Server struct {
	Router  *gin.Engine
	rootGrp *gin.RouterGroup
	cfg     *config.Config
}

func (s *Server) BindAPIRoutes(fn RegisterRoutesFunc, store storage.Storage) {
	fn(s.rootGrp, s.cfg, store)
}

func New(cfg *config.Config) *Server {
	s := &Server{
		Router: gin.New(),
		cfg:    cfg,
	}

	s.rootGrp = s.Router.Group("/v1")

	return s
}
