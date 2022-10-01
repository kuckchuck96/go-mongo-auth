package controller

import (
	"go-mongo-auth/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMessage(ctx *gin.Context) {
	var message service.Message
	message.New()

	if err := ctx.ShouldBindJSON(&message); err != nil {
		log.Println("Unable to bind request body.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	service.CreateMessage(message)
}
