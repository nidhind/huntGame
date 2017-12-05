package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendJSON(c *gin.Context, b interface{}, e error) {
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
	} else {
		c.JSON(http.StatusOK, &b)
	}
}
