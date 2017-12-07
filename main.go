package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nidhind/huntGame/db"
)

var startUpTime time.Time = time.Now()

func main() {

	// Initialize database
	db.InitMongo()

	api := gin.New()
	// Prevent redirects on trailing slashes
	api.RedirectTrailingSlash = false
	// Enable Logger
	api.Use(gin.Logger())
	// Enable CROS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	api.Use(cors.New(corsConfig))

	// Mount API routes
	mountRoutes(api)

	// Default port is 8080
	// To override set PORT env variable
	api.Run()
}
