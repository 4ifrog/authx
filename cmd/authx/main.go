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
	store := mongo.New(cfg)
	defer store.Close()
	if err := store.SeedUserData(); err != nil {
		panic(err)
	}

	// HTTP server
	srv := server.New(cfg)
	srv.BindAPIRoutes(api.GetRoutesFunc(), store)

	// Put everything in an app and run it.
	a := app.New(srv, cfg)
	a.Run()
}
