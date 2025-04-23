package controllers

import (
	"Cars/internal/models"
	"Cars/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Роуттарды тіркеу
func RegisterCarRoutes(router *gin.Engine) {
	router.GET("/cars", getCars)
	router.POST("/cars", createCar)
	router.PUT("/cars/:id", updateCar)
	router.DELETE("/cars/:id", deleteCar)
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
	services.DeleteCar(id)
	c.JSON(http.StatusOK, gin.H{"message": "Car deleted"})
}
