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
		Engine *gin.Engine
		Jwt    jwt.IJwtToken
	}
)

func NewMiddleware(engine *gin.Engine, jwt jwt.IJwtToken) IMiddleware {
	return &Middleware{
		Engine: engine,
		Jwt:    jwt,
	}
}

func (m *Middleware) AddMiddlewares() {
	middlewares := []gin.HandlerFunc{
		requestValidation(m),
	}

	m.Engine.Use(middlewares...)
}
