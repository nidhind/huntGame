// handlers
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func statusHandler(c *gin.Context) {
	c.JSON(200, map[string](interface{}){
		"name":   "huntGame",
		"env":    os.Getenv("GO_ENV"),
		"uptime": fmt.Sprint(time.Since(startUpTime)),
	})
}
