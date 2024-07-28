package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"net/http"
	"os"
	"time"
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

	usera, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Printf("Error in CompleteUserAuth: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User
	result := database.DB.Where("email = ?", usera.Email).First(&dbUser)

	if result.Error != nil {
		fmt.Printf("User not found in database, creating new user\n")
		// If user does not exist, create a new user
		dbUser = models.User{
			Username: usera.NickName,
			Email:    usera.Email,
			UserType: models.TypeMember,
			IsActive: true,
		}
		if err := database.DB.Create(&dbUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		fmt.Printf("User found in database: %+v\n", dbUser)
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": dbUser.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the token as an HTTP-only cookie
	c.SetCookie(
		"token",
		tokenString,
		3600*24*30, // 30 days
		"/",
		"https://note-taking-dusky.vercel.app", // Change this to your domain
		true,
		true,
	)

	fmt.Println("Authentication successful, sending response")
	c.Redirect(http.StatusFound, os.Getenv("CALLBACK_URL"))
}
