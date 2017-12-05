package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
	"github.com/nidhind/huntGame/utils"

	"golang.org/x/crypto/bcrypt"
)

func addUserHandler(c *gin.Context) {

	emailId := c.Param("emailId")
	if !utils.IsValidEmail(emailId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
			"status":  "failed",
			"code":    "1000",
			"message": "Invalid emailId",
		})
		return
	}

	// Parse request body into JSON
	var user models.SignUpReq
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
			"status":  "failed",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}

	if !utils.IsValidPassword(user.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
			"status":  "failed",
			"code":    "1001",
			"message": "Invalid or weak password",
		})
		return
	}

	// Check if user already exists
	if DoesEmailExists(emailId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
			"status":  "failed",
			"code":    "1003",
			"message": "User already exists",
		})
		return
	}

	// Hash the password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	user.Password = string(hash)

	utils.SendJSON(c, &user, nil)
}

// Check if user exists by emailId
func DoesEmailExists(id string) bool {
	_, exists := GetUserByEmail(id)
	return exists
}

// fetch user from DB by email
func GetUserByEmail(id string) (db.User, bool) {
	var user db.User

	return user, false
}
