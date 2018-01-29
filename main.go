package main

import (
	"time"
	"os"
	"io"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/utils"
)

var startUpTime time.Time = time.Now()

func main() {

	// Initialize database
	db.InitMongo()

	// Initialize audit logger
	utils.InitAuditLog()

	api := gin.New()

	// Disable Console Color when writing the logs to file.
	gin.DisableConsoleColor()
	// Logging to a file.
	f, _ := os.OpenFile("gin.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	gin.DefaultWriter = io.MultiWriter(f,os.Stdout)

	// Prevent redirects on trailing slashes
	api.RedirectTrailingSlash = false
	// Enable Logger
	api.Use(gin.Logger())
	// Enable CROS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowMethods("PATCH")
	api.Use(cors.New(corsConfig))

	// Mount API routes
	mountRoutes(api)

	// Default port is 8080
	// To override set PORT env variable
	api.Run()
}
