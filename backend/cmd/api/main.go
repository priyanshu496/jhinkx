package main

import (
	"fmt"
	"log"

	"github.com/priyanshu496/jhinkx.git/internal/config"
	"github.com/priyanshu496/jhinkx.git/internal/db" // Import our new db package

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load the configuration from our .env file
	cfg := config.Load()

	// 2. Initialize the Database connection
	// We pass the DatabaseURL we loaded from the .env file into our new InitDB function
	db.InitDB(cfg.DatabaseURL)

	// 3. Initialize the Gin router
	router := gin.Default()

	// 4. Set up a health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Server is healthy and connected to DB!",
		})
	})

	// 5. Start the server
	fmt.Printf("Starting Teambuilder API on port: %s\n", cfg.Port)
	
	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}