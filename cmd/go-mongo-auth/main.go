package main

import (
	"fmt"
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"
	"go-mongo-auth/internal/middleware"
	"go-mongo-auth/internal/route"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

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

	if err := engine.Run(fmt.Sprintf(":%v", config.App.Port)); err != nil {
		log.Fatalln(err)
	}
}
