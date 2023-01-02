package route

import (
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/controller"
	"go-mongo-auth/internal/pkg/jwt"
	"go-mongo-auth/internal/pkg/mongo"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type (
	IRoute interface {
		AddRoutes()
	}

	route struct {
		engine    *gin.Engine
		appConfig config.App
		user      controller.IUserController
		health    controller.IHealth
	}
)

func NewRoute(engine *gin.Engine, config config.Config, jwt jwt.IJwtToken, mongoClient mongo.IMongoClient) IRoute {
	return &route{
		engine:    engine,
		appConfig: config.App,
		user:      controller.NewUserController(config, jwt, mongoClient),
		health:    controller.NewHealthController(),
	}
}

func (r *route) AddRoutes() {
	v1 := r.engine.Group(r.appConfig.BasePath)
	{
		userRoutes(r, v1)
		healthRoutes(r, v1)
		swaggerRoutes(r, v1)
	}
}

func userRoutes(r *route, group *gin.RouterGroup) {
	userRoutes := group.Group("/user")
	userRoutes.POST("/login", r.user.Login)
	userRoutes.POST("/register", r.user.Register)
}

func healthRoutes(r *route, group *gin.RouterGroup) {
	userRoutes := group.Group("/health")
	userRoutes.GET("/liveness", r.health.HealthCheck)
}

func swaggerRoutes(r *route, group *gin.RouterGroup) {
	swaggerRoutes := group.Group("/swagger")
	swaggerRoutes.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
