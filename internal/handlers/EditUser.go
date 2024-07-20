package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func EditUser(c *gin.Context) {
	// Get the ID of the user making the request
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get the ID of the user to be edited
	idToEdit := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, idToEdit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the current user is an admin
	var currentUser models.User
	if err := database.DB.First(&currentUser, currentUserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current user"})
		return
	}
	isAdmin := currentUser.IsAdmin

	var updateData struct {
		Username  string          `json:"username"`
		Email     string          `json:"email"`
		FirstName string          `json:"first_name"`
		LastName  string          `json:"last_name"`
		Password  string          `json:"password"`
		UserType  models.UserType `json:"user_type"`
		IsActive  *bool           `json:"is_active"`
		IsAdmin   *bool           `json:"is_admin"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user is trying to edit their own profile or if they're an admin
	if currentUserID != user.ID && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this user"})
		return
	}

	// Fields that both normal users and admins can update
	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(updateData.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Fields that only admins can update
	if isAdmin {
		if updateData.Username != "" {
			user.Username = updateData.Username
		}
		if updateData.UserType != "" {
			user.UserType = updateData.UserType
		}
		if updateData.IsActive != nil {
			user.IsActive = *updateData.IsActive
		}
		if updateData.IsAdmin != nil {
			user.IsAdmin = *updateData.IsAdmin
		}
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Remove sensitive information
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}
