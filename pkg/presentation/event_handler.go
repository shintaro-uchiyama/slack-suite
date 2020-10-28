package presentation

import (
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventHandler struct{}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

type EventCreateRequest interface {}

func (h EventHandler) Create(c *gin.Context) {
	var eventCreateRequest EventCreateRequest
	if err := c.ShouldBind(&eventCreateRequest); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}

	logrus.Info(eventCreateRequest)
	c.JSON(http.StatusOK, nil)
}
