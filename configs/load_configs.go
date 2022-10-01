package configs

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var v *viper.Viper

const configs = "configs"

func init() {
	v = viper.New()
}

func Load(args string) {
	profile := getProfile(args)
	if profile == "" {
		profile = configs
	} else {
		profile = strings.Join([]string{configs, profile}, "-")
	}

	log.Printf("Loading '%v' profile.\n", profile)

	v.AddConfigPath("configs/")
	v.SetConfigName(profile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Panicln("Unable to load config.", err)
	}
}

func Get(prop string) string {
	return v.GetString(prop)
}

func GetChrono(prop string) time.Duration {
	return v.GetDuration(prop)
}

func getProfile(args string) string {
	if strings.TrimSpace(args) == "" {
		return ""
	}

	profile := strings.Split(args, "=")
	if profile[0] == "profile" {
		return profile[1]
	}
	return ""
}
