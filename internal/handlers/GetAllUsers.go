package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAllUsers retrieves all users
func GetAllUsers(c *gin.Context) {
    // Get the ID of the user making the request
    currentUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Check if the current user is an admin
    if !currentUserIsAdmin(currentUserID) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view all users"})
        return
    }

    var users []models.User

    // You might want to add pagination here for large datasets
    if err := database.DB.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }

    // Remove sensitive information
    for i := range users {
        users[i].Password = ""
    }

    c.JSON(http.StatusOK, gin.H{"users": users})
}

