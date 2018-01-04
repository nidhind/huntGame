package main

import (
	"net/http"
  "fmt"
	"github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
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

  level := puzzle.Level
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
    Level:       level,
    Image:       puzzle.Image,
    Clue:        puzzle.Clue,
    SolutionHash: puzzle.SolutionHash,
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
