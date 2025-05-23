package controllers

import (
	"Cars/internal/middleware"
	"Cars/internal/models"
	"Cars/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SuperAdminController struct {
	service *services.SuperAdminService
}

func NewSuperAdminController(service *services.SuperAdminService) *SuperAdminController {
	return &SuperAdminController{service: service}
}

// RegisterSuperAdminRoutes registers all superadmin routes
func RegisterSuperAdminRoutes(router *gin.Engine, controller *SuperAdminController) {
	superadmin := router.Group("/api/superadmin")
	superadmin.Use(middleware.AuthMiddleware(), RequireSuperAdmin())

	superadmin.GET("/users", controller.GetAllUsers)
	superadmin.GET("/users/:id", controller.GetUser)
	superadmin.PUT("/users/:id/role", controller.UpdateUserRole)
	superadmin.PUT("/users/:id/block", controller.BlockUser)
	superadmin.PUT("/users/:id/unblock", controller.UnblockUser)
}

// RequireSuperAdmin middleware to check if user is a super admin
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if userModel, ok := user.(*models.User); !ok || userModel.Role != string(models.RoleSuperAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: requires super admin role"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetAllUsers returns all users except super admins
func (c *SuperAdminController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUser returns a specific user
func (c *SuperAdminController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateUserRole updates a user's role
func (c *SuperAdminController) UpdateUserRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var request struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := c.service.UpdateUserRole(uint(id), request.Role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user role updated successfully"})
}

// BlockUser blocks a user
func (c *SuperAdminController) BlockUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := c.service.BlockUser(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user blocked successfully"})
}

// UnblockUser unblocks a user
func (c *SuperAdminController) UnblockUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := c.service.UnblockUser(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user unblocked successfully"})
}
