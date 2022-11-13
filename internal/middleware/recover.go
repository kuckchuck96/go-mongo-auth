package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultResponse struct {
	Message     string
	Status      int
	Description string
}

func AppRecovery(ctx *gin.Context, recovered any) {
	res := DefaultResponse{
		Message: http.StatusText(http.StatusInternalServerError),
		Status:  http.StatusInternalServerError,
	}
	if err, ok := recovered.(string); ok {
		res.Description = err
	}
	ctx.AbortWithStatusJSON(res.Status, res)
}
