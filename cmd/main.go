package main

import (
	"benzinga-webhook-receiver/internal/api"
	"benzinga-webhook-receiver/pkg/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig("pkg/config/config.yaml")
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting the webhook receiver application")

	router := gin.Default()

	api.SetupRoutes(router)
	config.StartBatchTimer() // Start the batch processing timer

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
