// Handels most of the user specific operations

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nidhind/huntGame/db"
	"github.com/nidhind/huntGame/models"
	"github.com/nidhind/huntGame/utils"

	"time"
	"golang.org/x/crypto/bcrypt"
)

func addUserHandler(c *gin.Context) {

	// Parse request body into JSON
	var user models.SignUpReq
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}
	emailId := user.EmailId
	if !utils.IsValidEmail(emailId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1000",
			"message": "Invalid emailId",
		})
		return
	}

	if !utils.IsValidPassword(user.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1001",
			"message": "Invalid or weak password",
		})
		return
	}

	// Check if user already exists
	if DoesEmailExists(emailId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1003",
			"message": "User already exists",
		})
		return
	}

	// Hash the password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	// Create User Object to insert
	u := db.InsertUserQuery{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Password:    user.Password,
		Email:       emailId,
		AccessLevel: "normal",
		Level:       1,
		PreviousLevelFinishTime: time.Now(),
	}
	err = db.InsertNewUser(&u)
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

// Fetch and serve user profile
func getUserProfile(c *gin.Context) {
	// This is a authenticated route
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
	r := models.ProfileRes{
		Code:   "0",
		Status: "success",
		Payload: &models.UserProfile{
			FirstName:               u.FirstName,
			LastName:                u.LastName,
			Email:                   u.Email,
			Level:                   u.Level,
			LevelImage:              p.Image,
			LevelClue:               p.Clue,
			AccessLevel:             u.AccessLevel,
			AccessToken:             u.AccessToken,
			PreviousLevelFinishTime: u.PreviousLevelFinishTime.String(),
		},
	}
	c.JSON(http.StatusOK, &r)
}

// Fetch and serve user leader board
func getUserLeadBoardHandler(c *gin.Context) {
	l := 10
	ul, err := db.GetUserLeaderBoard(l)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
			"status":  "error",
			"code":    "500",
			"message": "Internal server error",
		})
		return
	}

	payload := &[]models.UserLeaderBoard{}
	for _, u := range *ul {
		*payload = append(*payload, models.UserLeaderBoard{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Level:     u.Level,
			PreviousLevelFinishTime: u.PreviousLevelFinishTime.String(),
		})
	}

	r := models.ProfileRes{
		Code:    "0",
		Status:  "success",
		Payload: &payload,
	}
	c.JSON(http.StatusOK, &r)
}

//udpate user roles
func roleHandler(c *gin.Context) {

	var role models.RoleUpdateReq
	err := c.ShouldBindJSON(&role)
	if err != nil {
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

	//validate user role
	if u.AccessLevel == "admin" {
		if DoesEmailExists(role.Email){
			err = db.UpdateRoleByEmailId(role.Email,role.AccessLevel)
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
				"message": "user role updated",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1010",
			"message": "User does not exist",
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, &map[string](interface{}){
		"status":  "error",
		"code":    "1009",
		"message": "action requires higher access level",
	})
}
