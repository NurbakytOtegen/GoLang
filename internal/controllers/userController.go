package controllers

import (
	"github.com/Kabyl/Cars-Goland-Project/internal/models"
	"github.com/Kabyl/Cars-Goland-Project/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine) {
	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
}

func getUsers(c *gin.Context) {
	users := services.GetUsers()
	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	services.CreateUser(user)
	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	services.UpdateUser(id, user)
	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	services.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
