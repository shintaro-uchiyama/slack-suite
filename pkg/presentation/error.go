package presentation

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) *errorResponse {
	return &errorResponse{
		Message: err.Error(),
	}
}

func jsonError(c *gin.Context, code int, err error) {
	logrus.Error(err.Error())
	c.JSON(code, NewErrorResponse(err))
}
