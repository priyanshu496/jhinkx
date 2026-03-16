package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a person looking for a space on the platform.
type User struct {
	ID                 string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username           string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash       string         `gorm:"not null" json:"-"`
	PreferredGroupSize int            `gorm:"not null" json:"preferred_group_size"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"` // Enterprise feature: Soft deletes
}

// Space (formerly Group) represents a chat room where matched users are placed.
type Space struct {
	ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TargetSize int            `gorm:"not null" json:"target_size"`
	Status     string         `gorm:"type:varchar(20);default:'active'" json:"status"` // E.g., 'active', 'deleting'
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// SpaceMember is the junction table linking Users to the Spaces they are in.
type SpaceMember struct {
	SpaceID       string    `gorm:"type:uuid;primaryKey" json:"space_id"`
	UserID        string    `gorm:"type:uuid;primaryKey" json:"user_id"`
	VotedToDelete bool      `gorm:"default:false" json:"voted_to_delete"` // NEW: Tracks the consensus delete votes!
	JoinedAt      time.Time `gorm:"autoCreateTime" json:"joined_at"`

	// SECURITY: Database-level foreign key constraints. 
	// If a space or user is deleted, their membership is safely cascaded (deleted) too.
	Space Space `gorm:"foreignKey:SpaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Message stores the chat history for a specific space.
type Message struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SpaceID   string         `gorm:"type:uuid;not null;index" json:"space_id"`
	UserID    string         `gorm:"type:uuid;not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Allows users to "unsend" messages in the future!

	// SECURITY: Link messages securely to their Space and User
	Space Space `gorm:"foreignKey:SpaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}