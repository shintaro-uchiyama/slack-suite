package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shintaro-uchiyama/pkg/infrastructure"

	"github.com/shintaro-uchiyama/pkg/application"

	"github.com/gin-gonic/gin"
	"github.com/shintaro-uchiyama/pkg/presentation"
	"github.com/sirupsen/logrus"
)

func initRoute() {
	secretManager, err := infrastructure.NewSecretManager()
	if err != nil {
		logrus.Fatal(fmt.Errorf("NewSecretManager error: %w", err))
	}
	eventSlack := infrastructure.NewEventSlack()
	slackApplication := application.NewVerifyApplication(secretManager, eventSlack)

	pubSub, err := infrastructure.NewPubSub()
	if err != nil {
		logrus.Fatal(fmt.Errorf("NewPubSub error: %w", err))
	}
	taskApplication := application.NewTaskApplication(pubSub)

	eventHandler := presentation.NewEventHandler(
		slackApplication,
		taskApplication,
	)

	r := gin.Default()
	r.Use(logMiddleWare())
	r.Use(errorMiddleWare())
	r.POST("/events", eventHandler.Create)
	err = r.Run()
	if err != nil {
		logrus.Fatal(fmt.Errorf("run gin engine error: %w", err))
	}
}

func logMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetOutput(os.Stdout)

		logger := logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"url":    c.Request.URL,
		})
		logger.Info("start")
		c.Next()
		logger.Info("end")
	}
}

func errorMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			logrus.Error(fmt.Errorf("public error: %w", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": err.Error(),
			})
			return
		}

		err = c.Errors.ByType(gin.ErrorTypePrivate).Last()
		if err != nil {
			logrus.Error(fmt.Errorf("private error: %w", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": "An unexpected error has occurred",
			})
			return
		}
	}
}

func main() {
	initRoute()
}
