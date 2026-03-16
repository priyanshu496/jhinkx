package db

import (
	"log"

	"github.com/priyanshu496/jhinkx.git/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(databaseURL string) {
	var err error

	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to the PostgreSQL database!")

	// --- NEW CODE: Auto Migration ---
	// This tells GORM: "Look at these structs. If tables for them don't exist, create them. 
	// If the structs have new fields, add those columns to the tables."
	err = DB.AutoMigrate(
		&models.User{},
		&models.Space{},
		&models.SpaceMember{},  
		&models.Message{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("Database tables migrated successfully!")
}