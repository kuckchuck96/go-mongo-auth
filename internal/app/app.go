package app

import (
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"
	"go-mongo-auth/internal/middleware"
	"go-mongo-auth/internal/route"
	"go-mongo-auth/internal/swagger"

	"github.com/gin-gonic/gin"
)

func Initialize(engine *gin.Engine, config config.Config) error {
	// Jwt token
	jwt := jwt.NewJwtToken(config)

	// Configure Middlewares
	middleware.NewMiddleware(engine, jwt).AddMiddlewares()

	// Init mongo
	mongo, err := database.NewMongoClient(config.Mongo)
	if err != nil {
		return err
	}

	// Configure routes
	route.NewRoute(engine, config, jwt, mongo).AddRoutes()

	// Configure swagger
	swagger.ConfigureSwagger(config.App)

	return nil
}
