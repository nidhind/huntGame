// handlers
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func statusHandler(c *gin.Context) {
	c.JSON(200, map[string](string){"name": "huntGame", "env": os.Getenv("GO_ENV"),
		"uptime": fmt.Sprint(time.Since(startUpTime))})
}

func addUserHandler(c *gin.Context) {
	body, err := addUser()
	if body != nil {
		c.JSON(http.S)
	}
}
