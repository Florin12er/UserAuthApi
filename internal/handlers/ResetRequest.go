package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"UserAuth/pkg/utils"
	"net/http"
	"time"
	"crypto/rand"
    "encoding/base32"
	"sync"
	"github.com/gin-gonic/gin"
)

var (
	resetCodes = make(map[string]resetCodeInfo)
	resetMutex sync.Mutex
)

type resetCodeInfo struct {
	Code      string
	ExpiresAt time.Time
}

func generateResetCode() (string, error) {
    // Generate 6 random bytes
    randomBytes := make([]byte, 6)
    _, err := rand.Read(randomBytes)
    if err != nil {
        return "", err
    }
    
    // Encode the bytes to base32 and return the first 6 characters
    return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)[:6], nil
}


func ResetRequest(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	resetCode , err:= generateResetCode()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset code"})
        return
    }
	expirationTime := time.Now().Add(15 * time.Minute)

	resetMutex.Lock()
	resetCodes[request.Email] = resetCodeInfo{Code: resetCode, ExpiresAt: expirationTime}
	resetMutex.Unlock()

	if err := utils.SendResetEmail(request.Email, resetCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset code sent to email"})
}

