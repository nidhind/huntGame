package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
	"github.com/nidhind/huntGame/utils"
  "fmt"
  "time"
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
  answerHash := utils.GenerateHash(answer)
  fmt.Println("answer is ",answer)
  fmt.Println("answerHash is ",answerHash)

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

  fmt.Println("Solution hash is",p.SolutionHash)

  if p.SolutionHash == answerHash {
    // update user profile level and leaderboard
    fmt.Println("Level ",u.Level, "finished at ", time.Now())
    fmt.Println(u.FirstName," has advanced to Level ",u.Level+1)
    fmt.Println("Solution hash is",p.SolutionHash)
    c.JSON(http.StatusOK, &map[string](interface{}){
      "code":    "0",
      "status":  "success",
      "message": "Correct answer",
    })
    return
  } else {
      fmt.Println(u.FirstName," has advanced to ",u.Level+1)
      fmt.Println("Level ",u.Level, "finished at ", time.Now())
      c.JSON(http.StatusOK, &map[string](interface{}){
      "code":    "1008",
      "status":  "failure",
      "message": "Incorrect answer",
    })
    return
  }

}
