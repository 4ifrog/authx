package config

import (
	"github.com/spf13/pflag"
)

// Mapstructure maps the fields to the program parameter names.
type Config struct {
	Port          int
	Debug         bool
	RedisAddr     string `mapstructure:"redis-addr"`
	MongoAddr     string `mapstructure:"mongo-addr"`
	AccessSecret  string `mapstructure:"access-secret"`
	RefreshSecret string `mapstructure:"refresh-secret"`
	AccessTTL     int    `mapstructure:"access-ttl"`
	RefreshTTL    int    `mapstructure:"refresh-ttl"`
}

func setDefaults(flagset *pflag.FlagSet) {
	flagset.Bool("debug", false, "Enable debug")
	flagset.Int32("port", 8080, "The TCP port to run the application")
	flagset.String("redis-addr", "localhost:6379", "The address of Redis to where the app connects")
	flagset.String("mongo-addr", "mongodb://nobody:secrets@localhost:27017/authx", "The address of Mongo to where the app connects")
	flagset.String("access-secret", "", "Secret for signing an access token")
	flagset.String("refresh-secret", "", "Secret for signing a refresh token")
	flagset.Int32("access-ttl", 24 * 60 * 60, "Access token TTL in seconds")
	flagset.Int32("refresh-ttl", 30 * 24 * 60 * 60, "Refresh token TTL in seconds")
}
