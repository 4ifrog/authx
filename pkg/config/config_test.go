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
redis-addr: redis.net:35380
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

func TestDefaultValues(t *testing.T) {
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
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
}

func TestLoadFromConfigFile(t *testing.T) {
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
	assert.Equal(t, "redis.net:35380", cfg.RedisAddr)
}

func TestLoadFromEnvVariables(t *testing.T) {
	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	Setenv(t, "CS_PORT", "6000")
	Setenv(t, "CS_DEBUG", "true")
	Setenv(t, "CS_REDIS_ADDR", "10.20.30.40")
	defer func() {
		Unsetenv(t, "CS_PORT")
		Unsetenv(t, "CS_DEBUG")
		Unsetenv(t, "CS_REDIS_ADDR")
	}()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args:    []string{},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 6000, cfg.Port)
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "10.20.30.40", cfg.RedisAddr)
}

func TestLoadFromFlags(t *testing.T) {
	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args: []string{
			"--debug", "true",
			"--port", "7000",
			"--redis-addr", "redis.example.com:30000",
		},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 7000, cfg.Port)
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "redis.example.com:30000", cfg.RedisAddr)
}

func TestOverrides(t *testing.T) {
	viper.Reset()

	cfg := New()
	v := viper.GetViper()

	writeYAML(t)
	defer func() {
		require.NoError(t, os.Remove("config.yaml"))
	}()

	Setenv(t, "CS_DEBUG", "false")
	defer func() {
		Unsetenv(t, "CS_DEBUG")
	}()

	opts := BindConfigOpts{
		FlagSet: getTestFlagSet(),
		Args: []string{
			"--redis-addr", "redis.example.com:30000",
		},
	}
	cfg.BindConfig(v, opts)
	cfg.LoadConfig(v)

	assert.Equal(t, 5000, cfg.Port)
	assert.Equal(t, false, cfg.Debug)
	assert.Equal(t, "redis.example.com:30000", cfg.RedisAddr)
}
