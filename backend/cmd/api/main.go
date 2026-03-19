package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/priyanshu496/jhinkx.git/internal/api"
	"github.com/priyanshu496/jhinkx.git/internal/config"
	"github.com/priyanshu496/jhinkx.git/internal/db" // Import our new db package
	"github.com/priyanshu496/jhinkx.git/internal/kafka"
	"github.com/priyanshu496/jhinkx.git/internal/redis"
	"github.com/priyanshu496/jhinkx.git/internal/ws"
)

func main() {
	// 1. Load the configuration from our .env file
	cfg := config.Load()

	// 2. Initialize the Database connection
	// We pass the DatabaseURL we loaded from the .env file into our new InitDB function
	db.InitDB(cfg.DatabaseURL)
	kafka.InitProducer(cfg)
	kafka.InitConsumer(cfg)
	redis.InitRedis()
	chatHub := ws.NewHub()
	go chatHub.Run()
	// 3. Initialize the Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow your Next.js app
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		// Space routes
		apiRoutes.GET("/spaces", api.GetUserSpaces)
		apiRoutes.GET("/spaces/:id", api.GetSpaceDetails)
		// Chat History and Consensus Delete
		apiRoutes.GET("/spaces/:id/messages", api.GetSpaceMessages)
		apiRoutes.POST("/spaces/:id/vote-delete", api.VoteDeleteSpace)
		apiRoutes.POST("/spaces/match", api.RequestMatch)
	}

	router.GET("/ws/spaces/:id", func(c *gin.Context) {
		ws.ServeWS(chatHub, c)
	})
	// 5. Start the server
	fmt.Printf("Starting Teambuilder API on port: %s\n", cfg.Port)

	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
