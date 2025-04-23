package controllers

import (
	"github.com/Kabyl/Cars-Goland-Project/internal/middleware"
	"github.com/Kabyl/Cars-Goland-Project/internal/models"
	"github.com/Kabyl/Cars-Goland-Project/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterAuthRoutes registers the authentication routes
func RegisterAuthRoutes(router *gin.Engine) {
	router.POST("/auth/register", register)
	router.POST("/auth/login", login)
	router.GET("/auth/me", middleware.AuthMiddleware(), getMe)
}

// register handles user registration
func register(c *gin.Context) {
	var registerRequest models.RegisterRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register user
	user, err := services.RegisterUser(registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	token, err := services.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return user and token
	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user.ToUserResponse(),
	})
}

// login handles user login
func login(c *gin.Context) {
	var loginRequest models.LoginRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login user
	user, token, err := services.LoginUser(loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return user and token
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user.ToUserResponse(),
	})
}

// getMe returns the current authenticated user
func getMe(c *gin.Context) {
	// Get user ID from context (set by authMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	// Get user from database
	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Return user
	c.JSON(http.StatusOK, user.ToUserResponse())
}
