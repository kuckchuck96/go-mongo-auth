package controller

import (
	"go-mongo-auth/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var login service.Login

	if err := ctx.ShouldBindJSON(&login); err != nil {
		log.Println("Unable to bind request body.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := service.Authenticate(login)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func Register(ctx *gin.Context) {
	var req service.User
	req.New()

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Unable to bind request body.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := service.Register(req)
	if err != nil {
		log.Println("Unable to create new user.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
