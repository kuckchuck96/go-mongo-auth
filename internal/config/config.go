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
		Uri      string
		Database string
		Timeout  time.Duration
	}

	Auth struct {
		Expiry time.Duration
	}

	Jwt struct {
		SigningKey string `mapstructure:"signing-key"`
		Auth       Auth
	}

	App struct {
		Name string
		Env  string
		Port int
	}
	Config struct {
		App   App
		Jwt   Jwt
		Mongo Mongo
	}
)

func NewConfig() (Config, error) {
	profile := flag.String("profile", "", "Profile for run configuration")
	flag.Parse()

	if strings.TrimSpace(*profile) == "" {
		return Config{}, errors.New("no profile specified for run configuration")
	}

	v := viper.New()
	v.AddConfigPath(configDir)
	v.SetConfigName(*profile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return Config{}, errors.New("unable to load configuration")
	}

	log.Printf("'%v' profile activated.\n", *profile)

	var c Config
	_ = v.Unmarshal(&c)

	return c, nil
}

// func Load() error {
// 	profile := flag.String("profile", "", "Profile for run configuration")
// 	flag.Parse()

// 	if strings.TrimSpace(*profile) == "" {
// 		return errors.New("no profile specified for run configuration")
// 	}

// 	v = viper.New()
// 	v.AddConfigPath(configDir)
// 	v.SetConfigName(*profile)
// 	v.SetConfigType("yaml")
// 	if err := v.ReadInConfig(); err != nil {
// 		return errors.New("unable to load configuration")
// 	}

// 	log.Printf("'%v' profile activated.\n", *profile)

// 	var c Config
// 	_ = v.Unmarshal(&c)
// 	log.Println("CONFIG:", c)

// 	return nil
// }

// func Get(prop string) string {
// 	return v.GetString(prop)
// }

// func GetChrono(prop string) time.Duration {
// 	return v.GetDuration(prop)
// }
