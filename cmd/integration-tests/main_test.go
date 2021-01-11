package tests

import (
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/cybersamx/codefresh/pkg/api"
	"github.com/cybersamx/codefresh/pkg/app"
	"github.com/cybersamx/codefresh/pkg/config"
	"github.com/cybersamx/codefresh/pkg/rdb"
	"github.com/cybersamx/codefresh/pkg/server"
)

var a *app.App

func TestMain(m *testing.M) {
	// TestMain needs to use os.Exit to wrap the bootstrapper function so that
	// the `defer` can be executed properly.
	os.Exit(bootstrap(m))
}

func bootstrap(m *testing.M) int {
	// Config
	v := viper.GetViper()
	cfg := config.New()
	cfg.BindConfig(v)
	cfg.LoadConfig(v)

	// Redis
	rc := rdb.New(cfg)
	defer rc.Close()

	// HTTP server
	srv := server.New(cfg)
	srv.BindAPIRoutes(api.GetRoutesFunc(), rc)

	// Put everything in an app and run it.
	a = app.New(srv, cfg)

	return m.Run()
}
