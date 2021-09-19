package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/cybersamx/authx/pkg/app"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/store/mongo"
)

var testapp *app.App

func TestMain(m *testing.M) {
	bootstrap()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func bootstrap() {
	fmt.Println("Bootstrap...")

	// Config
	v := viper.GetViper()
	cfg := config.New()
	cfg.BindConfig(v)
	cfg.LoadConfig(v)
	cfg.TemplatesDir = "web/templates"

	// Mongo
	ds := mongo.New(cfg)
	if err := ds.SeedUserData(); err != nil {
		panic(err)
	}

	// HTTP server
	srv := server.New(cfg)
	srv.BindAPIRoutes(GetRoutesFunc(), ds)

	// Put everything in an app and run it.
	testapp = app.New(srv, ds, cfg)
}

func teardown() {
	fmt.Println("Teardown...")

	testapp.Store.Close()
}
