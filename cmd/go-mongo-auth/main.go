package main

import (
	"fmt"
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"
	"go-mongo-auth/internal/middleware"
	"go-mongo-auth/internal/route"
	"go-mongo-auth/internal/swagger"
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

	// Jwt token
	jwt := jwt.NewJwtToken(config)

	// Configure Middlewares
	middleware.NewMiddleware(engine, jwt).AddMiddlewares()

	// Init mongo
	mongo, err := database.NewMongoClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	// Configure routes
	route.NewRoute(engine, config, jwt, mongo).AddRoutes()

	// Configure swagger
	swagger.ConfigureSwagger(config.App)

	if err := engine.Run(fmt.Sprintf(":%v", config.App.Port)); err != nil {
		log.Fatalln(err)
	}
}
