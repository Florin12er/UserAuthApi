package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"net/http"
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GithubAuth() {
	githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	gothic.Store = store
	goth.UseProviders(
		github.New(
			githubClientId,
			githubClientSecret,
			"http://localhost:8080/auth/github/callback",
		),
	)
}

func CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	gothic.GetProviderName = func(*http.Request) (string, error) {
		return provider, nil
	}

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Printf("Error in CompleteUserAuth: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User
	result := database.DB.Where("email = ?", user.Email).First(&dbUser)

	if result.Error != nil {
		fmt.Printf("User not found in database, creating new user\n")
		// If user does not exist, create a new user
		dbUser = models.User{
			Username: user.NickName,
			Email:    user.Email,
			UserType: models.TypeMember,
			IsActive: true,
		}
		database.DB.Create(&dbUser)
	} else {
		fmt.Printf("User found in database: %+v\n", dbUser)
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": dbUser.Username,
		"email":    dbUser.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set JWT as a cookie
	c.SetCookie("token", tokenString, 86400*3, "/", "localhost", false, true)

	fmt.Println("Authentication successful, sending response")

	c.JSON(http.StatusOK, gin.H{
		"message": "User authenticated successfully",
		"token":   tokenString,
	})
	c.Redirect(http.StatusFound, "/")
}
