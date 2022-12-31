package config

import (
	"errors"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const configDir = "./config"

type (
	Mongo struct {
		Uri            string
		Database       string
		Timeout        time.Duration
		ContextTimeout time.Duration `mapstructure:"context-timeout"`
	}

	Auth struct {
		Expiry time.Duration
	}

	Jwt struct {
		SigningKey string `mapstructure:"signing-key"`
		Auth       Auth
	}

	App struct {
		Name        string
		Env         string
		Description string
		BasePath    string `mapstructure:"base-path"`
		Version     string
	}

	Server struct {
		Port     int
		WaitTime time.Duration `mapstructure:"wait-time"`
	}

	Config struct {
		App    App
		Server Server
		Jwt    Jwt
		Mongo  Mongo
	}
)

const (
	_yaml    = "yaml"
	_profile = "profile"
)

func NewConfig() (Config, error) {
	var env string
	var c Config
	profile := flag.String(_profile, "", "Profile for run configuration")
	flag.Parse()
	if !flag.Parsed() {
		return c, errors.New("unable to parse the flag")
	}
	env = *profile

	if strings.TrimSpace(env) == "" {
		return c, errors.New("no profile specified for run configuration")
	}

	v := viper.New()
	v.AddConfigPath(configDir)
	v.SetConfigName(env)
	v.SetConfigType(_yaml)
	if err := v.ReadInConfig(); err != nil {
		return c, errors.New("unable to load configuration")
	}

	// Set env vars
	if err := setEnvVars(v); err != nil {
		return c, err
	}

	log.Printf("'%v' profile activated.\n", env)

	if err := v.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}

func setEnvVars(v *viper.Viper) error {
	if v == nil {
		return errors.New("viper instance is nil")
	}

	var envName string
	for _, key := range v.AllKeys() {
		envName = v.GetString(key)
		if strings.HasPrefix(envName, "$") {
			if err := v.BindEnv(key, strings.TrimPrefix(envName, "$")); err != nil {
				return err
			}
		}
	}

	return nil
}
