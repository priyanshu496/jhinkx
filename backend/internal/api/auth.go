package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/priyanshu496/jhinkx.git/internal/config"
	"github.com/priyanshu496/jhinkx.git/internal/db"
	"github.com/priyanshu496/jhinkx.git/internal/models"
)

// SignupInput defines the exact JSON structure we expect from the frontend
type SignupInput struct {
	Username           string `json:"username" binding:"required"`
	Password           string `json:"password" binding:"required,min=6"`
	PreferredGroupSize int    `json:"preferred_group_size" binding:"required,min=3"`
}

// SigninInput defines the JSON for logging in
type SigninInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Signup handles new account creation
func Signup(c *gin.Context) {
	var input SignupInput

	// 1. Validate the incoming JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hash the password securely
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Prepare the User model
	user := models.User{
		Username:           input.Username,
		PasswordHash:       string(hashedPassword),
		PreferredGroupSize: input.PreferredGroupSize,
	}

	// 4. Save to the database
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username might be taken."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signup successful!"})
}

// Signin handles user login and returns a JWT
func Signin(c *gin.Context) {
	var input SigninInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// 1. Find the user by username
	if err := db.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 2. Compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 3. Generate the JWT token containing the user's ID
	cfg := config.Load()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Expires in 72 hours
	})

	// 4. Sign the token
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// 5. Send token to frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Signin successful!",
		"token":   tokenString,
	})
}