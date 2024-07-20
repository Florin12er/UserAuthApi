package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"UserAuth/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/url" // Import the net/url package
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var registerData struct {
		Username          string `json:"username" binding:"required"`
		Email             string `json:"email" binding:"required,email"`
		Password          string `json:"password" binding:"required,min=8"`
		TurnstileResponse string `json:"cf-turnstile-response" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify Cloudflare Turnstile
	if !verifyTurnstile(registerData.TurnstileResponse) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CAPTCHA"})
		return
	}

	// Check if username or email already exists
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", registerData.Username, registerData.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
		return
	}

	// Validate password strength
	if !isValidPassword(registerData.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not meet the security requirements"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	user := models.User{
		Username: registerData.Username,
		Email:    registerData.Email,
		Password: string(hashedPassword),
		UserType: models.TypeMember,
		IsActive: false, // Set to false until email is verified
	}

	// Generate verification code
	verificationCode := utils.GenerateRandomCode(6) // Implement this function to generate a random 6-digit code
	user.VerificationCode = verificationCode

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send verification email
	if err := utils.SendVerificationEmail(user.Email, verificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful. Please check your email for verification."})
}

func verifyTurnstile(response string) bool {
	secretKey := os.Getenv("CLOUDFLARE_TURNSTILE_SECRET_KEY")
	verifyURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	// Check if we're in a test environment
	if os.Getenv("APP_ENV") == "test" && response == "mock_turnstile_response" {
		log.Println("Test mode: Bypassing Turnstile verification")
		return true
	}

	resp, err := http.PostForm(verifyURL, url.Values{
		"secret":   {secretKey},
		"response": {response},
	})
	if err != nil {
		log.Printf("Error making POST request to Turnstile verification endpoint: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Decode the response from Cloudflare
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding response from Turnstile verification endpoint: %v", err)
		return false
	}

	// Log the entire response for debugging purposes
	log.Printf("Turnstile verification response: %v", result)

	// Check if the verification was successful
	success, ok := result["success"].(bool)
	if !ok || !success {
		log.Printf("Turnstile verification failed: %v", result)
		return false
	}

	return true
}

// isValidPassword checks if the password meets the security requirements
func isValidPassword(password string) bool {
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

