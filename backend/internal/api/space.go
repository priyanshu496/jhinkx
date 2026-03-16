package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/priyanshu496/jhinkx.git/internal/db"
	"github.com/priyanshu496/jhinkx.git/internal/models"
)

// GetUserSpaces fetches all spaces the logged-in user is a member of
func GetUserSpaces(c *gin.Context) {
	// 1. Get the secure userID from the AuthMiddleware
	userID, _ := c.Get("userID")

	// 2. Query the junction table (SpaceMembers) to find their spaces
	var spaceMembers []models.SpaceMember
	
	// We use GORM's Preload("Space") to automatically grab the actual Space details 
	// along with the membership record in a single database query!
	if err := db.DB.Preload("Space").Where("user_id = ?", userID).Find(&spaceMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spaces"})
		return
	}

	// 3. Extract just the Space data to send a clean list to the frontend
	var spaces []models.Space
	for _, sm := range spaceMembers {
		spaces = append(spaces, sm.Space)
	}

	c.JSON(http.StatusOK, gin.H{"spaces": spaces})
}

// GetSpaceDetails fetches the metadata and members of a specific space
func GetSpaceDetails(c *gin.Context) {
	// 1. Grab the ID from the URL (e.g., /api/spaces/1234-5678)
	spaceID := c.Param("id")
	
	// Grab the logged-in user's ID
	userID, _ := c.Get("userID")

	// 2. SECURITY CHECK: Ensure this user is actually a member of this space!
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&membership).Error; err != nil {
		// If no record is found, they are either an attacker guessing IDs, or the space doesn't exist.
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this space"})
		return
	}

	// 3. Fetch the space details
	var space models.Space
	if err := db.DB.Where("id = ?", spaceID).First(&space).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Space not found"})
		return
	}

	// 4. Fetch the other members of this space (Preloading the User data to get usernames)
	var members []models.SpaceMember
	db.DB.Preload("User").Where("space_id = ?", spaceID).Find(&members)

	// 5. Format the members list so we don't accidentally send sensitive data (like password hashes)
	var memberProfiles []map[string]interface{}
	for _, m := range members {
		memberProfiles = append(memberProfiles, map[string]interface{}{
			"id":              m.User.ID,
			"username":        m.User.Username,
			"joined_at":       m.JoinedAt,
			"voted_to_delete": m.VotedToDelete,
		})
	}

	// 6. Return everything nicely to the frontend
	c.JSON(http.StatusOK, gin.H{
		"space":   space,
		"members": memberProfiles,
	})
}

// GetSpaceMessages fetches the chat history for a specific space
func GetSpaceMessages(c *gin.Context) {
	spaceID := c.Param("id")
	userID, _ := c.Get("userID")

	// 1. SECURITY CHECK: Is the user a member of this space?
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&membership).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to these messages"})
		return
	}

	// 2. Fetch the messages, ordered by the time they were sent (oldest first)
	var messages []models.Message
	if err := db.DB.Where("space_id = ?", spaceID).Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// VoteDeleteSpace handles the consensus deletion logic
func VoteDeleteSpace(c *gin.Context) {
	spaceID := c.Param("id")
	userID, _ := c.Get("userID")

	// 1. SECURITY CHECK: Is the user a member?
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&membership).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot vote to delete a space you are not in"})
		return
	}

	// 2. Record this user's vote to delete
	if err := db.DB.Model(&membership).Update("voted_to_delete", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register your vote"})
		return
	}

	// 3. THE CONSENSUS LOGIC: Check if everyone has voted!
	var totalMembers int64
	var deleteVotes int64

	// Count total members in the space
	db.DB.Model(&models.SpaceMember{}).Where("space_id = ?", spaceID).Count(&totalMembers)
	
	// Count how many of those members have voted "yes"
	db.DB.Model(&models.SpaceMember{}).Where("space_id = ? AND voted_to_delete = ?", spaceID, true).Count(&deleteVotes)

	// 4. If the votes match the total members, delete the space!
	if totalMembers > 0 && totalMembers == deleteVotes {
		// Because we set up our database models securely, deleting the space 
		// will also safely handle the space_members and messages attached to it!
		db.DB.Where("id = ?", spaceID).Delete(&models.Space{})
		
		c.JSON(http.StatusOK, gin.H{
			"message": "Consensus reached! The space has been permanently deleted.",
			"deleted": true,
		})
		return
	}

	// If we are still waiting on others to vote:
	c.JSON(http.StatusOK, gin.H{
		"message":      "Your vote has been recorded. Waiting for other members to agree.",
		"total_votes":  deleteVotes,
		"members_left": totalMembers - deleteVotes,
		"deleted":      false,
	})
}