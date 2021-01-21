package main

import (
	"github.com/spf13/viper"

	"github.com/cybersamx/authx/pkg/api"
	"github.com/cybersamx/authx/pkg/app"
	"github.com/cybersamx/authx/pkg/config"
	"github.com/cybersamx/authx/pkg/server"
	"github.com/cybersamx/authx/pkg/storage/mongo"
)

func main() {
	// Config
	v := viper.GetViper()
	cfg := config.New()
	cfg.BindConfig(v)
	cfg.LoadConfig(v)

	// Mongo
	mStore := mongo.New(cfg)
	defer mStore.Close()
	if err := mStore.SeedUserData(); err != nil {
		panic(err)
	}

	// HTTP server
	srv := server.New(cfg)
	srv.BindAPIRoutes(api.GetRoutesFunc(), mStore)

	// Put everything in an app and run it.
	a := app.New(srv, cfg)
	a.Run()
}
