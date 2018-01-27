package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func mountRoutes(app *gin.Engine) {

	// Get server status
	app.GET("/status", statusHandler)
	// Get user profile
	app.GET("/users/profile", authenticateToken, getUserProfile)
	// Get user leader board
	app.GET("/users/leader-board", authenticateToken, getUserLeadBoardHandler)
	// Password reset redirect request handler
	app.GET("/users/password/reset/new", ForgotPasswordRedirectHandler)

	// Add new user
	app.POST("/users", addUserHandler)
	//user login
	app.POST("/users/access_token", sessionHandler)
	//add new puzzle
	app.POST("/puzzle", authenticateToken, addPuzzleHandler)
	//submit answer
	app.POST("/answer", authenticateToken, answerHandler)
	//user role update
	app.POST("/users/role", authenticateToken, roleHandler)

	// Forgot password send mail handler
	app.PATCH("/users/password/reset", forgotPasswordHandler)
	// Password reset update handler
	app.PUT("/users/password/reset/new", ForgotPasswordUpdateHandler)

	//serve static assets
	// app.Static("/assets", "./assets")

	// Handle 404
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string](string){
			"message": "Resource not found",
		})
	})
}
