package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	var verifyData struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&verifyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", verifyData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.VerificationCode != verifyData.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Check if the code has expired (15 minutes)
	if time.Now().Sub(user.VerificationCodeSentAt) > 15*time.Minute {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification code has expired"})
		return
	}

	// Activate the user
	user.IsActive = true
	user.VerificationCode = ""
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

