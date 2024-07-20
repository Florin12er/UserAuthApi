package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetUserById retrieves a user by their ID
// GetUserById retrieves a user by their ID
func GetUserById(c *gin.Context) {
    // Get the ID of the user making the request
    currentUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Get the ID of the user to be retrieved
    idToGet := c.Param("id")

    var user models.User
    if err := database.DB.First(&user, idToGet).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Check if the current user is an admin or the same user
    if currentUserID != user.ID && !currentUserIsAdmin(currentUserID) {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this user"})
        return
    }

    // Remove sensitive information
    user.Password = ""

    c.JSON(http.StatusOK, gin.H{"user": user})
}

// Helper function to check if the current user is an admin
func currentUserIsAdmin(userID interface{}) bool {
    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        return false
    }
    return user.IsAdmin
}

