package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Port          int
	Debug         bool
	RedisAddr     string `mapstructure:"redis-addr"`
	AccessSecret  string `mapstructure:"access-secret"`
	RefreshSecret string `mapstructure:"refresh-secret"`
}

type BindConfigOpts struct {
	FlagSet *pflag.FlagSet
	Args    []string
}

// Config load order: is: default values > config file > environment variables > CLI arguments
// For details on load precedence see https://github.com/spf13/viper#why-viper
func (c *Config) BindConfig(v *viper.Viper, set ...BindConfigOpts) {
	// Config flag
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".") // Working directory

	// CLI flags
	if len(set) == 0 {
		set = []BindConfigOpts{
			{
				FlagSet: pflag.CommandLine,
				Args:    os.Args,
			},
		}
	}
	opts := set[0]

	opts.FlagSet.Bool("debug", false, "Enable debug")
	opts.FlagSet.Int32("port", 8080, "The TCP port to run the application")
	opts.FlagSet.String("redis-addr", "localhost:6379", "The address of Redis to where the application connects")
	opts.FlagSet.String("access-secret", "", "Secret for signing an access token")
	opts.FlagSet.String("refresh-secret", "", "Secret for signing a refresh token")

	if err := opts.FlagSet.Parse(opts.Args); err != nil {
		panic(fmt.Errorf("failed to parse arguments: %v", err))
	}
	if err := v.BindPFlags(opts.FlagSet); err != nil {
		panic(fmt.Errorf("failed to bind pflags: %v", err))
	}

	// Environment variables
	// The setup allows the following mappings of env vars and flags (key)
	// CS_PORT <--> port
	v.AutomaticEnv()
	v.SetEnvPrefix("CS")
	replacer := strings.NewReplacer("-", "_")
	v.SetEnvKeyReplacer(replacer)
}

func (c *Config) LoadConfig(v *viper.Viper) {
	// If there's no config file to load, it's ok and move on.
	_ = v.ReadInConfig()

	if err := v.UnmarshalExact(c); err != nil {
		panic(fmt.Errorf("failed to umarshal parsed configurations to config: %v", err))
	}
}

func New() *Config {
	return &Config{}
}
