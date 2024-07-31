package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CheckAuth verifies the JWT token from the cookie and checks if the user is authenticated.
func CheckAuth(c *gin.Context) {
	// Get the cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"authenticated": false, "message": "No authentication token found"},
		)
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
		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				c.JSON(
					http.StatusUnauthorized,
					gin.H{"authenticated": false, "message": "Token has expired"},
				)
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "message": "Invalid expiration claim"})
			return
		}

		// Token is valid
		userId, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"authenticated": false, "message": "Invalid user ID claim"},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"authenticated": true, "user_id": uint(userId)})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "message": "Invalid token claims"})
	}
}
