package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router.POST("/user", handlers.SendEmail)

	router.Run()
}
