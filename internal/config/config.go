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

var v *viper.Viper

func Load() error {
	profile := flag.String("profile", "", "Profile for run configuration")
	flag.Parse()

	if strings.TrimSpace(*profile) == "" {
		return errors.New("no profile specified for run configuration")
	}

	v = viper.New()
	v.AddConfigPath(configDir)
	v.SetConfigName(*profile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return errors.New("unable to load configuration")
	}

	log.Printf("'%v' profile activated.\n", *profile)

	return nil
}

func Get(prop string) string {
	return v.GetString(prop)
}

func GetChrono(prop string) time.Duration {
	return v.GetDuration(prop)
}
