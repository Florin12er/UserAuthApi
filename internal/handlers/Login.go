package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxLoginAttempts = 5
	lockoutDuration  = 15 * time.Minute
)

func Login(c *gin.Context) {
	var loginUser struct {
		EmailOrUsername string `json:"email_or_username" binding:"required"`
		Password        string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Try to find the user by email or username
	if err := database.DB.Where("email = ? OR username = ?", loginUser.EmailOrUsername, loginUser.EmailOrUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if the account is locked
	if time.Now().Before(user.LockedUntil) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is locked. Please try again later."})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)); err != nil {
		// Increment failed login attempts
		user.FailedLoginAttempts++
		user.LastFailedLogin = time.Now()

		if user.FailedLoginAttempts >= maxLoginAttempts {
			user.LockedUntil = time.Now().Add(lockoutDuration)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account locked due to too many failed attempts. Please try again later."})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}

		database.DB.Save(&user)
		return
	}

	// Reset failed login attempts on successful login
	user.FailedLoginAttempts = 0
	user.LockedUntil = time.Time{}
	database.DB.Save(&user)

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the token as an HTTP-only cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30, // 30 days
		Path:     "/",
		Domain:   "note-taking-dusky.vercel.app",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": user.Username})
}

func ProtectedRoute(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected route!"})
}
