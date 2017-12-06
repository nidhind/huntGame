// Handels most of the user specific operations

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
			"status":  "error",
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
			"status":  "error",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}

	if !utils.IsValidPassword(user.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
			"status":  "error",
			"code":    "1001",
			"message": "Invalid or weak password",
		})
		return
	}

	// Check if user already exists
	if DoesEmailExists(emailId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
			"status":  "error",
			"code":    "1003",
			"message": "User already exists",
		})
		return
	}

	// Hash the password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	user.Password = string(hash)
	// Create User Object to insert
	u := db.InsertUserQuery{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Password:    user.Password,
		Email:       emailId,
		AccessLevel: "normal",
		Rank:        0,
	}
	err = db.InsertNewUser(&u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string](interface{}){
			"status":  "error",
			"code":    "500",
			"message": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusCreated, map[string](interface{}){
		"status":  "success",
		"code":    "0",
		"message": "created",
	})
}

// Check if user exists by emailId
func DoesEmailExists(id string) bool {
	_, err := db.GetUserByEmail(id)
	if err != nil && err.Error() == "not found" {
		// User doesnot exists
		return false
	} else if err != nil {
		panic(err)
	}
	// User exists
	return true
}
