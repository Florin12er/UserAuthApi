// internal/handlers/auth_handlers.go

package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckAuth(c *gin.Context) {
	// If the middleware.CheckAuthenticated() passes, it means the user is authenticated
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}
