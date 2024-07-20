package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteUser deletes a user by their ID
func DeleteUser(c *gin.Context) {
    // Get the ID of the user making the request
    currentUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Check if the current user is an admin
    if !currentUserIsAdmin(currentUserID) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete users"})
        return
    }

    // Get the ID of the user to be deleted
    idToDelete := c.Param("id")

    var user models.User
    if err := database.DB.First(&user, idToDelete).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if err := database.DB.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

