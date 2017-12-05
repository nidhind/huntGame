package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

func MountRoutes(app *gin.Engine) {

	// Enable CROS
	app.Use(cors.Default())

	// Get server status
	app.GET("/status", statusHandler)

	// Add new user
	app.POST("/users", addUserHandler)
	
	// Handle 404
	app.Use(func(c *gin.Context) {
		c.String(404, "Resource not found")
	})

}
