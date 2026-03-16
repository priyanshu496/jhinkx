package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// Using your exact module name with the .git extension
	"github.com/priyanshu496/jhinkx.git/internal/db"
	"github.com/priyanshu496/jhinkx.git/internal/models"
)

// GetCurrentUser fetches the profile of the currently logged-in user
func GetCurrentUser(c *gin.Context) {
	// 1. Grab the userID that our AuthMiddleware safely injected into the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	// 2. Look up the user in PostgreSQL
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3. Return the user data 
	// (Remember: The password hash is safely hidden because of the json:"-" tag in our struct!)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdatePreferencesInput defines what data a user is allowed to change
type UpdatePreferencesInput struct {
	PreferredGroupSize int `json:"preferred_group_size" binding:"required,min=3"`
}

// UpdatePreferences allows a user to change their settings
func UpdatePreferences(c *gin.Context) {
	// Grab the ID of the logged-in user making the request
	userID, _ := c.Get("userID")

	var input UpdatePreferencesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user's preferred group size in the database
	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("preferred_group_size", input.PreferredGroupSize).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated successfully!"})
}