package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

var startUpTime time.Time = time.Now()

func main() {

	api := gin.New()

	// Mount API routes
	MountRoutes(api)

	// Default port is 8080
	// To override set PORT env variable
	api.Run()
}
