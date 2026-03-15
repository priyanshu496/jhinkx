package models

import (
	"time"
)

// User represents a person looking for a group on the platform.
type User struct {
	// gen_random_uuid() tells PostgreSQL to automatically generate a secure UUID for us
	ID                 string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username           string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash       string    `gorm:"not null" json:"-"` // The "-" json tag ensures passwords are never accidentally sent to the frontend!
	PreferredGroupSize int       `gorm:"not null" json:"preferred_group_size"`
	CreatedAt          time.Time `json:"created_at"`
}

// Group represents a chat room where matched users are placed.
type Group struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TargetSize int       `gorm:"not null" json:"target_size"`
	CreatedAt  time.Time `json:"created_at"`
}

// GroupMember is the junction table linking Users to the Groups they are in.
type GroupMember struct {
	GroupID  string    `gorm:"type:uuid;primaryKey" json:"group_id"`
	UserID   string    `gorm:"type:uuid;primaryKey" json:"user_id"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

// Message stores the chat history for a specific group.
type Message struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	GroupID   string    `gorm:"type:uuid;not null;index" json:"group_id"` // Index makes searching for a group's messages super fast
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}