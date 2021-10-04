package main

import (
	"github.com/spf13/viper"

	"github.com/cybersamx/authx/pkg/api"
	"github.com/cybersamx/authx/pkg/app"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/store/mongo"
)

func main() {
	// Config
	v := viper.GetViper()
	cfg := config.New()
	cfg.BindConfig(v)
	cfg.LoadConfig(v)

	// Mongo
	ds := mongo.New(cfg)
	defer ds.Close()
	if err := api.SeedUserData(ds); err != nil {
		panic(err)
	}

	// HTTP server
	srv := server.New(cfg)
	srv.BindAPIRoutes(api.GetRoutesFunc(), ds)

	// Put everything in an app and run it.
	a := app.New(srv, ds, cfg)
	a.Run()
}
