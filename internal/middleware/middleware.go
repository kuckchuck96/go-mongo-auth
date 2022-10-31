package middleware

import (
	"go-mongo-auth/internal/jwt"

	"github.com/gin-gonic/gin"
)

type (
	IMiddleware interface {
		AddMiddlewares()
	}

	Middleware struct {
		engine *gin.Engine
		jwt    jwt.IJwtToken
	}
)

func NewMiddleware(engine *gin.Engine, jwt jwt.IJwtToken) IMiddleware {
	return &Middleware{
		engine: engine,
		jwt:    jwt,
	}
}

func (m *Middleware) AddMiddlewares() {
	middlewares := []gin.HandlerFunc{
		requestValidation(m),
	}

	m.engine.Use(middlewares...)
}
