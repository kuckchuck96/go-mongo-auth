package route

import (
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/controller"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"

	"github.com/gin-gonic/gin"
)

type (
	IRoute interface {
		AddRoutes()
	}

	Route struct {
		Engine *gin.Engine
		User   controller.IUserController
	}
)

func NewRoute(engine *gin.Engine, config config.Config, jwt jwt.IJwtToken, mongoClient database.IMongoClient) IRoute {
	return &Route{
		Engine: engine,
		User:   controller.NewUserController(config, jwt, mongoClient),
	}
}

func (r *Route) AddRoutes() {
	v1 := r.Engine.Group("/api/v1")
	{
		userRoutes(r, v1)
	}
}

func userRoutes(r *Route, group *gin.RouterGroup) {
	userRoutes := group.Group("user")
	userRoutes.POST("/login", r.User.Login)
	userRoutes.POST("/register", r.User.Register)
}
