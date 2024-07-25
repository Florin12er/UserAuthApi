package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"UserAuth/pkg/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(c *gin.Context) {
	var request struct {
		Email       string `json:"email" binding:"required,email"`
		ResetCode   string `json:"reset_code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resetMutex.Lock()
	codeInfo, exists := resetCodes[request.Email]
	resetMutex.Unlock()

	if !exists || codeInfo.Code != request.ResetCode || time.Now().After(codeInfo.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset code"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	resetMutex.Lock()
	delete(resetCodes, request.Email)
	resetMutex.Unlock()
	go func() {
		if err := utils.SendResetNotification(user.Email); err != nil {
			fmt.Println("Failed to send email:", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
