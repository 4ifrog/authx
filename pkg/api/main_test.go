package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/cybersamx/authx/pkg/app"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/store/mongo"
	"github.com/spf13/viper"
)

var testapp *app.App

func TestMain(m *testing.M) {
	fmt.Println("Main app staring up...")
	var code int

	// Config
	v := viper.GetViper()
	cfg := config.New()
	cfg.BindConfig(v)
	cfg.LoadConfig(v)
	cfg.TemplatesDir = "web/templates"

	// Mongo
	ds := mongo.New(cfg)
	defer func() {
		fmt.Println("Main app tearing down...")
		ds.Close()
		os.Exit(code)
	}()

	if err := SeedUserData(ds); err != nil {
		panic(err)
	}

	// HTTP server
	srv := server.New(cfg)
	srv.BindAPIRoutes(GetRoutesFunc(), ds)

	// Put everything in an app and run it.
	testapp = app.New(srv, ds, cfg)

	code = m.Run()
}
