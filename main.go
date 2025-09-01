package main

import (
	"github.com/gin-gonic/gin"
	"github.com/narumayase/anysher/log"
	"github.com/narumayase/anysher/middleware"
	"github.com/narumayase/anysher/middleware/gateway"
)

func main() {
	// Set the global log level
	log.SetLogLevel()

	// Create a new Gin router
	router := gin.New()

	// Use the middlewares
	router.Use(middleware.RequestIDToLogger())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())
	router.Use(gateway.Sender())

	// Define a sample route
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "world"})
	})

	// Start the server
	router.Run(":8080")
}
