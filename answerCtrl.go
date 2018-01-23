package main

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
	"github.com/nidhind/huntGame/utils"
)

func answerHandler(c *gin.Context) {

	// Parse request body into JSON
	var answerReq models.AnswerReq
	err := c.ShouldBindJSON(&answerReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}

	answer := answerReq.Answer

	// This is an authenticated route
	// User will be already present in context
	i, _ := c.Get("user")
	u := i.(*db.User)
	p, err := db.GetPuzzleByLevel(u.Level)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
			"status":  "error",
			"code":    "500",
			"message": "Puzzle data missing",
		})
		return
	}

	currentHash := utils.GenerateHash(answer)
	if bytes.Compare(currentHash, p.SolutionHash) != 0 {
		c.JSON(http.StatusOK, &map[string](interface{}){
			"code":    "1008",
			"status":  "failure",
			"message": "Incorrect answer",
		})
		return
	} else {
		//update user level since correct answer
		err = db.UpdateLevelByEmailId(u.Email, u.Level+1)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
				"status":  "error",
				"code":    "500",
				"message": "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, &map[string](interface{}){
			"code":    "0",
			"status":  "success",
			"message": "Correct answer",
		})
		return
	}
}
