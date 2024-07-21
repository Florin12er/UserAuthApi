package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

var (
	googleOAuthConfig *oauth2.Config
	googleStateString = "random" // Replace with a secure random string in production
)

func init() {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	url := googleOAuthConfig.AuthCodeURL(googleStateString, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != googleStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code exchange failed"})
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	service, err := googleoauth2.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OAuth2 service"})
		return
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Check if user already exists
	var user models.User
	if err := database.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		// Create new user
		user = models.User{
			Username: userInfo.Name,
			Email:    userInfo.Email,
			UserType: models.TypeMember,
			IsActive: true,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Create JWT token
	tokenString, err := createToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie(
		"token",
		tokenString,
		3600*24*30, // 30 days
		"/",
		"localhost", // Change this to your domain
		false,       // Set to true if using HTTPS
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func createToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

