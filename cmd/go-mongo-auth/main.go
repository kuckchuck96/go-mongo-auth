package main

import (
	"fmt"
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/middleware"
	"go-mongo-auth/internal/route"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configs
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	engine := gin.Default()
	engine.Use(middleware.RequestValidation())

	// Init mongo
	if err := database.ConnectionManager(); err != nil {
		log.Fatalln(err)
	}

	// Configure routes
	route.AddRoutes(engine)

	engine.Run(fmt.Sprintf(":%v", config.Get("app.port")))
}
