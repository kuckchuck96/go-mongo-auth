package main

import (
	"fmt"
	"go-mongo-auth/internal/app"
	"go-mongo-auth/internal/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Set gin mode
	mode := gin.ReleaseMode
	if config.App.Env == "dev" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.Default()

	// Initialize
	if err = app.Initialize(engine, config); err != nil {
		log.Fatalln(err)
	}

	if err := engine.Run(fmt.Sprintf(":%v", config.App.Port)); err != nil {
		log.Fatalln(err)
	}
}
