package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		"",
		false, // Set to true if using HTTPS
		true,
	)
	// Clear the user cookie
	c.SetCookie(
		"user",
		"",
		-1,
		"/",
		"",
		false, // Set to true if using HTTPS
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
