// session management

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
	"github.com/nidhind/huntGame/utils"
)

func sessionHandler(c *gin.Context) {

	// Parse request body into JSON
	var user models.GenAccessToken
	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}

	//validate user
	if !AuthenticateUser(user.Email, user.Password) {

		c.AbortWithStatusJSON(http.StatusUnauthorized, &map[string](interface{}){
			"status":  "error",
			"code":    "1005",
			"message": "Email or password is invalid.",
		})
		return
	}
	resp := models.LoginRes{
		Status: "success",
		Code:   "0",
	}
	resp.Payload.AccessToken = utils.GenerateAccessToken()
	err = db.UpdateAccessTokenByEmailId(user.Email, resp.Payload.AccessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
			"status":  "error",
			"code":    "500",
			"message": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, &resp)
	return
}

func AuthenticateUser(id string, password string) bool {
	user, err := db.GetUserByEmail(id)
	if err != nil && err.Error() == "not found" {
		// User does not exist
		return false
	} else if err != nil {
		panic(err)
	}
	// User exists
	error := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if error != nil {
		return false
	}
	return true
}

// Authenticate token and fetch user
func authenticateToken(c *gin.Context) {
	t := c.GetHeader("Authorization")
	if t == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "2000",
			"message": "Authorization parameters are invalid.",
		})
		return
	}

	u, err := db.GetUserByAccessToken(t)

	if err != nil && err.Error() == "not found" {
		// User does not exist
		c.AbortWithStatusJSON(http.StatusUnauthorized, &map[string](interface{}){
			"status":  "error",
			"code":    "2001",
			"message": "Invalid access-token",
		})
		return
	} else if err != nil {
		panic(err)
	}
	c.Set("user", &u)
	c.Next()
	return
}
