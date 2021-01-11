package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
)

type App struct {
	config *config.Config

	Router *gin.Engine
}

func (a *App) Run() {
	log.Println("Running...")
	log.Println("config", a.config)

	addr := fmt.Sprintf(":%d", a.config.Port)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func New(s *server.Server, cfg *config.Config) *App {
	app := &App{
		Router: s.Router,
		config: cfg,
	}

	return app
}
