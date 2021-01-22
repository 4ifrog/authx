package tests

import (
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/cybersamx/authx/pkg/api"
	"github.com/cybersamx/authx/pkg/app"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/storage/mongo"
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

	// Mongo
	store := mongo.New(cfg)
	defer store.Close()
	if err := store.SeedUserData(); err != nil {
		panic(err)
	}

	// HTTP server
	srv := server.New(cfg)
	srv.BindAPIRoutes(api.GetRoutesFunc(), store)

	// Put everything in an app and run it.
	a = app.New(srv, cfg)

	return m.Run()
}
