package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)



type BindConfigOpts struct {
	FlagSet *pflag.FlagSet
	Args    []string
}

const (
	envConfigPrefix   = "AX"
)

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

	setDefaults(opts.FlagSet)

	if err := opts.FlagSet.Parse(opts.Args); err != nil {
		panic(fmt.Errorf("failed to parse arguments: %v", err))
	}
	if err := v.BindPFlags(opts.FlagSet); err != nil {
		panic(fmt.Errorf("failed to bind pflags: %v", err))
	}

	// Environment variables
	// The setup allows the following mappings of env vars and flags (key)
	// AX_PORT <--> port
	v.AutomaticEnv()
	v.SetEnvPrefix(envConfigPrefix)
	// Use _ in environment variables and - in program params.
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
