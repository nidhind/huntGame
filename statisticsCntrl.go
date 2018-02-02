// get various statistics
package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/nidhind/huntGame/db"
	"net/http"
)

func statisticsHandler(c *gin.Context) {

	answerCount, err := db.GetCount("answer")
	if err != nil && err.Error() == "not found" {	
		// Create statistics Object to insert
		stat := db.InsertStatisticsQuery{
			Id: "answer",
			Count:	0,
		}
		err = db.InsertStatistics(&stat)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
				"status":  "error",
				"code":    "500",
				"message": "Internal server error",
			})
			return
		}
	}
	c.JSON(200, map[string](interface{}){
		"count":   answerCount,
	})
}
