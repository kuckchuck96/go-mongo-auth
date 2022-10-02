package route

import (
	"go-mongo-auth/internal/controller"

	"github.com/gin-gonic/gin"
)

func AddRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	{
		userRoutes(v1)
	}
}

func userRoutes(group *gin.RouterGroup) {
	userRoutes := group.Group("user")
	userRoutes.POST("/login", controller.Login)
	userRoutes.POST("/register", controller.Register)
}
