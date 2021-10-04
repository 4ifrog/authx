package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appName = "cmd"
)

func getTestFlagSet() *pflag.FlagSet {
	return pflag.NewFlagSet(appName, pflag.ContinueOnError)
}

// CREDITS:
// Based on https://github.com/spf13/viper/blob/master/internal/testutil/env.go
// Licensed under the MIT license
// Copyright (c) 2014 Steve Francia

func Setenv(t *testing.T, name, val string) {
	setenv(t, name, val, true)
}

func Unsetenv(t *testing.T, name string) {
	setenv(t, name, "", false)
}

func Getenv(t *testing.T, name string) string {
	return os.Getenv(name)
}

func setenv(t *testing.T, name, val string, valOK bool) {
	oldVal, oldOK := os.LookupEnv(name)

	if valOK {
		require.NoError(t, os.Setenv(name, val))
	} else {
		require.NoError(t, os.Unsetenv(name))
	}

	t.Cleanup(func() {
		if oldOK {
			require.NoError(t, os.Setenv(name, oldVal))
		} else {
			require.NoError(t, os.Unsetenv(name))
		}
	})
}

func writeYAML(t *testing.T) {
	content := `
port: 5000
debug: true
mongo-addr: mongo.net:35380
`

	filename := "config.yaml"
	file, err := os.Create(filename)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, file.Close())
	}()

	_, err = file.Write([]byte(content))
	require.NoError(t, err)
}

func Test_DefaultValues(t *testing.T) {
	// Need to unset os environment variables so that we can test defaults properly as
	// environment variables takes precedence over defaults.
	port := Getenv(t, "AX_PORT")
	debug := Getenv(t, "AX_DEBUG")
	mongoAddr := Getenv(t, "AX_MONGO_ADDR")
	Unsetenv(t, "AX_PORT")
	Unsetenv(t, "AX_DEBUG")
	Unsetenv(t, "AX_MONGO_ADDR")
	defer func() {
		Setenv(t, "AX_PORT", port)
		Setenv(t, "AX_DEBUG", debug)
		Setenv(t, "AX_MONGO_ADDR", mongoAddr)
	}()

	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args:    []string{},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, false, cfg.Debug)
	assert.Equal(t, "mongodb://nobody:secrets@localhost:27017/authx", cfg.MongoAddr)
}

func Test_LoadFromConfigFile(t *testing.T) {
	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	writeYAML(t)
	defer func() {
		require.NoError(t, os.Remove("config.yaml"))
	}()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args:    []string{},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 5000, cfg.Port)
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "mongo.net:35380", cfg.MongoAddr)
}

func Test_LoadFromEnvVariables(t *testing.T) {
	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	Setenv(t, "AX_PORT", "6000")
	Setenv(t, "AX_DEBUG", "true")
	Setenv(t, "AX_MONGO_ADDR", "10.20.30.40")
	defer func() {
		Unsetenv(t, "AX_PORT")
		Unsetenv(t, "AX_DEBUG")
		Unsetenv(t, "AX_MONGO_ADDR")
	}()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args:    []string{},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 6000, cfg.Port)
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "10.20.30.40", cfg.MongoAddr)
}

func Test_LoadFromFlags(t *testing.T) {
	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args: []string{
			"--debug", "true",
			"--port", "7000",
			"--mongo-addr", "mongo.example.com:30000",
		},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 7000, cfg.Port)
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "mongo.example.com:30000", cfg.MongoAddr)
}

func Test_Overrides(t *testing.T) {
	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	// Need to unset os environment variables so that we can test defaults properly as
	// environment variables takes precedence over defaults.
	port := Getenv(t, "AX_PORT")
	mongoAddr := Getenv(t, "AX_MONGO_ADDR")
	defer func() {
		Setenv(t, "AX_PORT", port)
		Setenv(t, "AX_MONGO_ADDR", mongoAddr)
	}()

	writeYAML(t)
	defer func() {
		require.NoError(t, os.Remove("config.yaml"))
	}()

	Setenv(t, "AX_DEBUG", "false")
	defer func() {
		Unsetenv(t, "AX_DEBUG")
	}()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args: []string{
			"--mongo-addr", "mongo.example.com:30000",
		},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 5000, cfg.Port)
	assert.Equal(t, false, cfg.Debug)
	assert.Equal(t, "mongo.example.com:30000", cfg.MongoAddr)
}
