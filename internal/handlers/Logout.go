package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		"localhost", // Change this to your domain
		false,       // Set to true if using HTTPS
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

