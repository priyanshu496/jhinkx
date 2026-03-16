package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/priyanshu496/jhinkx.git/internal/config"
)

// AuthMiddleware intercepts requests to check for a valid JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Look for the "Authorization" header in the incoming request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort() // Stops the request from going any further
			return
		}

		// 2. The token should look like "Bearer eyJhbGciOiJIUzI1NiIs..."
		// We need to split that string to get just the token part.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. Parse and validate the token using our secret key
		cfg := config.Load()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Extract the data (claims) we hid inside the token when they signed in
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Save the user_id into the Gin context so the next function can use it!
			c.Set("userID", claims["user_id"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to process token claims"})
			c.Abort()
			return
		}

		// 5. If everything is good, let them pass to the actual API endpoint
		c.Next()
	}
}