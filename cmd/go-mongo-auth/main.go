package main

import (
	"fmt"
	"go-mongo-auth/configs"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/middleware"
	"go-mongo-auth/internal/route"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	var args string

	if len(os.Args) > 1 {
		args = os.Args[1]
	} else {
		args = ""
	}

	configs.Load(args)
}

func main() {
	engine := gin.Default()
	engine.Use(middleware.RequestValidation())

	// Init mongo
	database.ConnectionManager()

	// Configure routes
	route.AddRoutes(engine)

	engine.Run(fmt.Sprintf(":%v", configs.Get("app.port")))
}
