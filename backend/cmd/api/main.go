package main

import (
	"fmt"
	"log"

	"github.com/priyanshu496/jhinkx.git/internal/api"
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

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Hello World",
		})
	})

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", api.Signup)
		authRoutes.POST("/signin", api.Signin)
	}
	
	// --- PROTECTED ROUTES ---
	// We group everything under /api and apply the AuthMiddleware
	apiRoutes := router.Group("/api")
	apiRoutes.Use(api.AuthMiddleware()) 
	{
		// These routes will now verify the JWT token before running!
		apiRoutes.GET("/users/me", api.GetCurrentUser)
		apiRoutes.PUT("/users/settings", api.UpdatePreferences)
	}
	// 5. Start the server
	fmt.Printf("Starting Teambuilder API on port: %s\n", cfg.Port)

	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
