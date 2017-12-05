package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
)

var startUpTime time.Time = time.Now()

func main() {
	// Initialize database
	db.InitMongo()

	api := gin.New()

	// Mount API routes
	MountRoutes(api)

	// Default port is 8080
	// To override set PORT env variable
	api.Run()
}
