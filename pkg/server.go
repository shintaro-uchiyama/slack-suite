package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shintaro-uchiyama/pkg/presentation"
	"github.com/sirupsen/logrus"
)

func initRoute() {
	eventHandler := presentation.NewEventHandler()
	r := gin.Default()
	r.Use(logMiddleWare())
	r.POST("/events", eventHandler.Create)
	err := r.Run()
	if err != nil {
		logrus.Fatal(fmt.Errorf("[main.initRoute]: %w", err))
	}
}

func logMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"url":    c.Request.URL,
		})
		logger.Info("start")
		c.Next()
		logger.Info("end")
	}
}

func main() {
	initRoute()
}
