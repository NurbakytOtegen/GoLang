package controllers

import (
	"net/http"

	"Cars/internal/middleware"
	"Cars/internal/models"
	"Cars/internal/services"

	"github.com/gin-gonic/gin"
)

// // Роуттарды тіркеу
// func RegisterCarRoutes(router *gin.Engine) {
// 	router.GET("/cars", getCars)
// 	router.POST("/cars", createCar)
// 	router.PUT("/cars/:id", updateCar)
// 	router.DELETE("/cars/:id", deleteCar)
// }

func RegisterCarRoutes(router *gin.Engine) {
	router.GET("/cars", getCars)
	router.GET("/cars/:id", middleware.AuthMiddleware(), GetCarByID)

	carGroup := router.Group("/cars")
	carGroup.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))

	{
		carGroup.POST("", createCar)
		carGroup.PUT("/:id", updateCar)
		carGroup.DELETE("/:id", deleteCar)
	}
}

// Барлық көліктерді алу
func getCars(c *gin.Context) {
	cars := services.GetCars()
	c.JSON(http.StatusOK, cars)
}

// Көлік қосу
func createCar(c *gin.Context) {
	var car models.Car
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	services.CreateCar(car)
	c.JSON(http.StatusCreated, car)
}

// Көлікті жаңарту
func updateCar(c *gin.Context) {
	id := c.Param("id")
	var car models.Car
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	services.UpdateCar(id, car)
	c.JSON(http.StatusOK, car)
}

// Көлікті өшіру
func deleteCar(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteCar(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
}

func GetCarByID(c *gin.Context) {
	id := c.Param("id")
	car, err := services.GetCarByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}
	c.JSON(http.StatusOK, car)
}
