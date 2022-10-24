package controller

import (
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"
	"go-mongo-auth/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	IUserController interface {
		Login(*gin.Context)
		Register(*gin.Context)
	}

	UserController struct {
		Service service.IUserService
	}
)

func NewUserController(config config.Config, jwt jwt.IJwtToken, mongoClient database.IMongoClient) IUserController {
	return &UserController{
		service.NewUserService(config, jwt, mongoClient),
	}
}

func (c *UserController) Login(ctx *gin.Context) {
	var login service.Login

	if err := ctx.ShouldBindJSON(&login); err != nil {
		log.Println("Unable to bind request body.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := c.Service.Authenticate(login)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *UserController) Register(ctx *gin.Context) {
	var req service.User

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Unable to bind request body.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := c.Service.Register(req)
	if err != nil {
		log.Println("Unable to create new user.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
