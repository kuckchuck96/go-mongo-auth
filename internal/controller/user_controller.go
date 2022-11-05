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

	userController struct {
		userService service.IUserService
	}
)

func NewUserController(config config.Config, jwt jwt.IJwtToken, mongoClient database.IMongoClient) IUserController {
	return &userController{
		userService: service.NewUserService(config, jwt, mongoClient),
	}
}

// Login godoc
// @Summary User login
// @Schemes
// @Description User login via email and password
// @Tags login
// @Accept json
// @Produce json
// @Param req body service.Login true "User login request"
// @Success 200 {object} service.AuthenticatedResponse
// @Failure 500 {object} service.UserErrResponse
// @Router /user/login [post]
func (c *userController) Login(ctx *gin.Context) {
	var login service.Login

	if err := ctx.ShouldBindJSON(&login); err != nil {
		log.Println("Unable to bind request body.", err)
		ctx.JSON(http.StatusInternalServerError, service.UserErrResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := c.userService.Authenticate(login)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, service.UserErrResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// Register godoc
// @Summary User registration
// @Schemes
// @Description User registration
// @Tags register
// @Accept json
// @Produce json
// @Param req body service.User true "User registeration request"
// @Success 200 {object} service.RegisteredResponse
// @Failure 500 {object} service.UserErrResponse
// @Router /user/register [post]
func (c *userController) Register(ctx *gin.Context) {
	var req service.User

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Unable to bind request body.", err)
		ctx.JSON(http.StatusInternalServerError, service.UserErrResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := c.userService.Register(req)
	if err != nil {
		log.Println("Unable to create new user.", err)
		ctx.JSON(http.StatusInternalServerError, service.UserErrResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
