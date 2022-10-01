package main

import (
	"fmt"
	"go-mongo-auth/configs"
	"go-mongo-auth/internal/controller"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/middleware"
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

	database.ConnectionManager()

	group := engine.Group("/api/v1")
	group.POST("/login", controller.Login)
	group.POST("/register", controller.Register)

	engine.Run(fmt.Sprintf(":%v", configs.Get("app.port")))
}
