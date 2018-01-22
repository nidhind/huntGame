package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
  "fmt"
  "time"
	"golang.org/x/crypto/bcrypt"
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
  //TODO format answer - remove spaces, convert to lowercase

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

  // fmt.Println(answer)
  // fmt.Println(p.SolutionHash)
  error := bcrypt.CompareHashAndPassword([]byte(p.SolutionHash),[]byte(answer))
  // fmt.Println(error)
  if error != nil {
        c.JSON(http.StatusOK, &map[string](interface{}){
        "code":    "1008",
        "status":  "failure",
        "message": "Incorrect answer",
      })
    return
  } else {
    // update user profile level and leaderboard
    fmt.Println("Level ",u.Level, "finished at ", time.Now())
    fmt.Println(u.FirstName," has advanced to Level ",u.Level+1)
    c.JSON(http.StatusOK, &map[string](interface{}){
      "code":    "0",
      "status":  "success",
      "message": "Correct answer",
    })
    return
  }

}
