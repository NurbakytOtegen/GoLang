package controllers

import (
	"Cars/internal/middleware"
	"Cars/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	profileService *services.ProfileService
}

type UpdateNameRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func NewProfileController(profileService *services.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

// RegisterProfileRoutes registers the profile routes
func RegisterProfileRoutes(router *gin.Engine, profileController *ProfileController) {
	profileGroup := router.Group("/api/profile")
	{
		profileGroup.Use(middleware.AuthMiddleware()) // Require authentication for all profile routes
		profileGroup.GET("/me", profileController.GetMyProfile)
		profileGroup.GET("/id/:id", profileController.GetProfileByID)
		profileGroup.GET("/email/:email", profileController.GetProfileByEmail)
		profileGroup.GET("/name/:name", profileController.GetProfileByName)
		profileGroup.PUT("/update/name", profileController.UpdateName)
		profileGroup.PUT("/update/password", profileController.UpdatePassword)
	}
}

// GetMyProfile handles GET /api/profile/me
func (c *ProfileController) GetMyProfile(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	user, reviews, err := c.profileService.GetUserProfile(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"reviews": reviews,
	})
}

// GetProfileByID handles GET /api/profile/id/:id
func (c *ProfileController) GetProfileByID(ctx *gin.Context) {
	userID := ctx.GetUint("id")
	user, reviews, err := c.profileService.GetUserProfile(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"reviews": reviews,
	})
}

// GetProfileByEmail handles GET /api/profile/email/:email
func (c *ProfileController) GetProfileByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	user, reviews, err := c.profileService.GetUserProfileByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"reviews": reviews,
	})
}

// GetProfileByName handles GET /api/profile/name/:name
func (c *ProfileController) GetProfileByName(ctx *gin.Context) {
	name := ctx.Param("name")
	user, reviews, err := c.profileService.GetUserProfileByName(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"reviews": reviews,
	})
}

// UpdateName handles PUT /api/profile/update/name
func (c *ProfileController) UpdateName(ctx *gin.Context) {
	var req UpdateNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := ctx.GetUint("userID")
	user, err := c.profileService.UpdateUserName(userID, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdatePassword handles PUT /api/profile/update/password
func (c *ProfileController) UpdatePassword(ctx *gin.Context) {
	var req UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := ctx.GetUint("userID")
	err := c.profileService.UpdateUserPassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "incorrect old password" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный текущий пароль"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
