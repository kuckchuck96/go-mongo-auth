package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	IHealth interface {
		HealthCheck(*gin.Context)
	}

	health struct{}
)

func NewHealthController() IHealth {
	return &health{}
}

func (c *health) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
