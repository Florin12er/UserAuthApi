package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(c *gin.Context) {
	// Get the cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "message": "No authentication token found"})
		return
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "message": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token has expired
		if float64(claims["exp"].(float64)) < float64(time.Now().Unix()) {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "message": "Token has expired"})
			return
		}

		// Token is valid
		c.JSON(http.StatusOK, gin.H{"authenticated": true, "user_id": claims["sub"]})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "message": "Invalid token claims"})
	}
}

