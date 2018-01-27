package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
	"github.com/nidhind/huntGame/utils"
)

func addPuzzleHandler(c *gin.Context) {
	// Parse request body into JSON
	var puzzle models.PuzzleReq
	err := c.ShouldBindJSON(&puzzle)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}

	//Authenticated route - user already in context
	i, _ := c.Get("user")
	u := i.(*db.User)
	if u.AccessLevel != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &map[string](interface{}){
			"status":  "error",
			"code":    "1009",
			"message": "action requires higher access level",
		})
		return
	}

	level := puzzle.Level
	hash := utils.GenerateHash(puzzle.SolutionHash)
	// Check if user already exists
	if DoesLevelExists(level) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1006",
			"message": "level already exists",
		})
		return
	}

	p := db.InsertPuzzleQuery{
		Level:        level,
		Image:        puzzle.Image,
		Clue:         puzzle.Clue,
		SolutionHash: hash,
	}

	err = db.InsertNewPuzzle(&p)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
			"status":  "error",
			"code":    "500",
			"message": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusCreated, &map[string](interface{}){
		"status":  "success",
		"code":    "0",
		"message": "created",
	})

}

//check if level exists
func DoesLevelExists(level int) bool {
	_, err := db.GetPuzzleByLevel(level)
	if err != nil && err.Error() == "not found" {
		// User doesnot exists
		return false
	} else if err != nil {
		panic(err)
	}
	// User exists
	return true
}
