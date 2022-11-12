package main

import (
	"fmt"
	"go-mongo-auth/internal/app"
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/pkg/serve"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const _dev = "dev"

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Set gin mode
	mode := gin.ReleaseMode
	if config.App.Env == _dev {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.Default()

	// Initialize
	if err = app.Initialize(engine, config); err != nil {
		log.Fatalln(err)
	}

	// Create custom http server
	serve.ListenAndServe(&http.Server{
		Addr:    fmt.Sprintf(":%v", config.Server.Port),
		Handler: engine,
	}, config.Server.WaitTime)
}
