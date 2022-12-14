package middleware

import (
	"go-mongo-auth/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type (
	IMiddleware interface {
		AddMiddlewares()
	}

	middleware struct {
		engine *gin.Engine
		jwt    jwt.IJwtToken
	}
)

func NewMiddleware(engine *gin.Engine, jwt jwt.IJwtToken) IMiddleware {
	return &middleware{
		engine: engine,
		jwt:    jwt,
	}
}

func (m *middleware) AddMiddlewares() {
	middlewares := []gin.HandlerFunc{
		gin.Logger(),
		gin.CustomRecovery(appRecovery),
		requestValidation(m),
	}

	m.engine.Use(middlewares...)
}
