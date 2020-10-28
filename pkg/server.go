package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shintaro-uchiyama/pkg/presentation"
	"github.com/sirupsen/logrus"
	"os"
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
		log := logrus.New()
		log.Level = logrus.DebugLevel
		log.Formatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "severity",
				logrus.FieldKeyMsg:   "message",
			},
		}
		log.Out = os.Stdout

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
