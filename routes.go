package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"fmt"

	"github.com/gin-contrib/cors"
)

func MountRoutes(app *gin.Engine) {

	// Enable CROS
	app.Use(cors.Default())

	// Get server status
	app.GET("/status", func(c *gin.Context) {
		c.JSON(200, map[string](string){"name": "huntGame", "env": os.Getenv("GO_ENV"),
			"uptime": fmt.Sprint(time.Since(startUpTime))})
	})

	// Handle 404
	app.Use(func(c *gin.Context) {
		c.String(404, "No middleware responded!")
	})

}
