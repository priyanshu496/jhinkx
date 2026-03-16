package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds all the environment variables we need for the app.
// Grouping them in a struct makes it easy to pass them around later.
type AppConfig struct {
	Port         string
	DatabaseURL  string
	KafkaBroker  string
	RedisURL     string
	JWTSecret   string
}

// Load reads the .env file and grabs the variables from the system.
func Load() *AppConfig {
	// Attempt to load the .env file.
	// We don't crash if it's missing, because in production (like on a cloud server),
	// environment variables are often set directly in the system, not via a .env file.
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found. Reading directly from system environment variables.")
	}

	// Create and return our configuration struct
	return &AppConfig{
		// os.Getenv grabs the value of the specific key from the environment
		Port:        getEnv("PORT", "8080"), // Default to 8080 if PORT isn't set
		DatabaseURL: os.Getenv("DATABASE_URL"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		RedisURL:    os.Getenv("REDIS_URL"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
	}
}

// getEnv is a small helper function. It tries to get an environment variable,
// but if it's empty, it returns a fallback value that we provide.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}