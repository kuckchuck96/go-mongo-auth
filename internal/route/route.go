package route

import (
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/controller"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type (
	IRoute interface {
		AddRoutes()
	}

	Route struct {
		engine    *gin.Engine
		appConfig config.App
		user      controller.IUserController
	}
)

func NewRoute(engine *gin.Engine, config config.Config, jwt jwt.IJwtToken, mongoClient database.IMongoClient) IRoute {
	return &Route{
		engine:    engine,
		appConfig: config.App,
		user:      controller.NewUserController(config, jwt, mongoClient),
	}
}

func (r *Route) AddRoutes() {
	v1 := r.engine.Group(r.appConfig.BasePath)
	{
		userRoutes(r, v1)
		swaggerRoutes(r, v1)
	}
}

func userRoutes(r *Route, group *gin.RouterGroup) {
	userRoutes := group.Group("/user")
	userRoutes.POST("/login", r.user.Login)
	userRoutes.POST("/register", r.user.Register)
}

func swaggerRoutes(r *Route, group *gin.RouterGroup) {
	swaggerRoutes := group.Group("/swagger")
	swaggerRoutes.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
