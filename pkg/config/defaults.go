package config

import (
	"github.com/spf13/pflag"
)

// Mapstructure maps the fields to the program parameter names.
type Config struct {
	Port               int
	Debug              bool
	MongoAddr          string `mapstructure:"mongo-addr"`
	AccessSecret       string `mapstructure:"access-secret"`
	RefreshSecret      string `mapstructure:"refresh-secret"`
	SessionSecret      string `mapstructure:"session-secret"`
	AccessTTL          int    `mapstructure:"access-ttl"`
	RefreshTTL         int    `mapstructure:"refresh-ttl"`
	RefreshTokenRotate bool   `mapstructure:"refresh-token-rotate"`
	StaticWebDir       string `mapstructure:"static-web-dir"`
	TemplatesDir       string `mapstructure:"templates-dir"`
}

func setDefaults(flagset *pflag.FlagSet) {
	flagset.Bool("debug", false, "Enable debug")
	flagset.Int32("port", 8080, "The TCP port to run the application") //nolint:gomnd
	flagset.String("mongo-addr", "mongodb://nobody:secrets@localhost:27017/authx", "The address of Mongo to where the app connects")
	flagset.String("access-secret", "", "Secret for signing an access token")
	flagset.String("refresh-secret", "", "Secret for signing a refresh token")
	flagset.Int32("access-ttl", 86400, "Access token TTL in seconds")    //nolint:gomnd
	flagset.Int32("refresh-ttl", 604800, "Refresh token TTL in seconds") //nolint:gomnd
	flagset.Bool("refresh-token-rotate", false, "Issue a new refresh token when renewing an access token")
	flagset.String("static-web-dir", "static", "The directory path containing the static web assets.")
	flagset.String("templates-dir", "templates", "The directory path containing the templates.")
}
