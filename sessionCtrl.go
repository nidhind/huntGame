// session management

package main

import (
	"net/http"
  "fmt"
	"github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
	"github.com/nidhind/huntGame/utils"
  "crypto/rand"
	"golang.org/x/crypto/bcrypt"
  "encoding/json"
  "log"
)

func sessionHandler(c *gin.Context) {


  // Parse request body into JSON
  var user models.LoginReq
  err := c.ShouldBindJSON(&user)
  fmt.Println(user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
      "status":  "error",
      "code":    "1002",
      "message": "Error in parsing JSON input",
    })
    return
  }

  // Check if user exist
  if !DoesEmailExists(user.Email) {
    c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
      "status":  "error",
      "code":    "1004",
      "message": "User does not exist",
    })
    return
  }

  // validate password
  if !utils.IsValidPassword(user.Password) {
    c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
      "status":  "error",
      "code":    "1005",
      "message": "Email or password is invalid.",
    })
    return
  }

  //validate user
  flag := validateUser(user.Email,user.Password)

  if  flag {
    accessToken := getAccessToken()
    fmt.Println(accessToken)

    var resp models.LoginRes
    resp.Status = "success"
    resp.Code = "0"
    resp.Message = "user validated"
    resp.Payload.AccessToken = accessToken

    out, err := json.Marshal(resp)
    if err != nil {
      log.Fatalln(err)
    }

    log.Println(string(out))
    c.JSON(http.StatusCreated,out)

    // c.JSON(http.StatusCreated, map[string](interface{}){
    //   "status":  "success",
    //   "code":    "0",
    //   "message": "user validated.",
    //   "payload" : {"access_token":string(accessToken),},
    // })

  } else {
    c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
      "status":  "error",
      "code":    "1005",
      "message": "Email or password is invalid.",
    })
    return
  }

}

func validateUser(id string, password string) bool {
  user, err := db.GetUserByEmail(id)
	if err != nil && err.Error() == "not found" {
		// User does not exist
		return false
	} else if err != nil {
		  panic(err)
	}
	// User exists
  error := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
  if error != nil {
    return false
  }
	return true
}

//generate access token
func getAccessToken() string {
	b := make([]byte, 64)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
