// Handels most of the user specific operations

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

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
			"code":    "5000",
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

// Send password reset mail
func forgotPasswordHandler(c *gin.Context) {
	// Parse request body into JSON
	var user models.ResetPswdEmailReq
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}
	email := user.Email
	if !utils.IsValidEmail(email) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1000",
			"message": "Invalid emailId",
		})
		return
	}

	// Check if user already exists
	if !DoesEmailExists(email) {
		// For security do not disclose that emailid does not exist
		time.Sleep(100 * time.Millisecond)
		c.AbortWithStatusJSON(http.StatusAccepted, &map[string](interface{}){
			"status":  "success",
			"code":    "0",
			"message": "Reset link has been emailed",
		})
		return
	}
	claims := map[string]interface{}{
		"email": email,
		"iat":   time.Now().Unix(),
		// Expires in 15 Minutes
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"sub": "RESET_PASSWORD_TOKEN",
	}
	rt := utils.GenerateJWTToken(claims)
	rl := GenerateResetPswdLink(&rt)

	// Send email with link
	pt, err := getResetPswdEmailTemplate()
	var msg bytes.Buffer
	td := struct {
		Email     string
		ResetLink string
	}{
		Email:     email,
		ResetLink: rl.String(),
	}
	err = pt.Execute(&msg, td)
	m := MailObj{
		To:          email,
		Subject:     `Reset your password`,
		MessageHTML: msg.String(),
	}
	err = Send(&m)
	if err != nil {
		log.Println("Error in sending mail: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
			"status":  "error",
			"code":    "500",
			"message": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusAccepted, &map[string](interface{}){
		"status":  "success",
		"code":    "0",
		"message": "Reset link has been emailed",
	})
}

func GenerateResetPswdLink(rt *string) url.URL {
	link, err := url.Parse(os.Getenv("RESET_PASSWORD_URL"))
	if err != nil {
		fmt.Println(err)
	}
	q := link.Query()
	q.Add("reset_token", *rt)
	link.RawQuery = q.Encode()
	return *link
}

// Validate reset token and update password
func ForgotPasswordRedirectHandler(c *gin.Context) {
	rt := c.Query("reset_token")
	_, err := utils.ParseJWTToken(rt)
	if err != "" {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusForbidden,
			`<h3>The link is either invalid or expired...</h3>`)
		return
	}
	rurl, perr := url.Parse(os.Getenv("RESET_PASSWORD_UPDATE_REDIRECT"))
	if perr != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusForbidden,
			`<h3>Unexpected error occured. Please contact admins...</h3>`)
		return
	}
	q := rurl.Query()
	q.Set("reset_token", rt)
	rurl.RawQuery = q.Encode()
	c.Redirect(http.StatusTemporaryRedirect, rurl.String())
}

func ForgotPasswordUpdateHandler(c *gin.Context) {
	jwt := c.Query("reset_token")
	claims, isValid := utils.ParseJWTToken(jwt)
	if isValid != "" {
		c.JSON(http.StatusForbidden, &map[string]interface{}{
			"status":  "error",
			"code":    "2003",
			"message": "Invalid or expired reset token",
		})
		return
	}

	// Parse request body into JSON
	var rp models.ResetPswdUpdateReq
	err := c.ShouldBindJSON(&rp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1002",
			"message": "Error in parsing JSON input",
		})
		return
	}

	// validate new password
	if !utils.IsValidPassword(rp.NewPassword) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &map[string](interface{}){
			"status":  "error",
			"code":    "1001",
			"message": "Invalid or weak password",
		})
		return
	}
	// TODO Handle unexpected errors during type casting
	e := claims["email"].(string)
	// Hash the password
	hash, _ := bcrypt.GenerateFromPassword([]byte(rp.NewPassword), 10)
	err = db.UpdatePasswordByEmailId(e, string(hash))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &map[string](interface{}){
			"status":  "error",
			"code":    "500",
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, &map[string](interface{}){
		"status":  "success",
		"code":    "0",
		"message": "New password applied",
	})
}
