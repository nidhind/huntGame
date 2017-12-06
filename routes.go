package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

func MountRoutes(app *gin.Engine) {

	// Prevent redirects on trailing slashes
	app.RedirectTrailingSlash = false

	// Enable CROS
	app.Use(gin.Logger())
	// Enable logger
	app.Use(cors.Default())

	// Get server status
	app.GET("/status", statusHandler)

	// Add new user
	app.POST("/users/:emailId", addUserHandler)

	// Handle 404
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string](string){
			"message": "Resource not found",
		})
	})
}
