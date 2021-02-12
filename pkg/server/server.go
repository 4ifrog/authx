package server

import (
	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/store"
)

type RegisterRoutesFunc func(router *gin.Engine, cfg *config.Config, store store.DataStore)

type Server struct {
	Router *gin.Engine
	cfg    *config.Config
}

func (s *Server) BindAPIRoutes(fn RegisterRoutesFunc, ds store.DataStore) {
	fn(s.Router, s.cfg, ds)
}

func New(cfg *config.Config) *Server {
	gin.SetMode(gin.ReleaseMode)
	s := &Server{
		Router: gin.Default(),
		cfg:    cfg,
	}

	return s
}
