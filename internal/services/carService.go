package services

import (
	"Cars/internal/models"
)

func GetCars() []models.Car {
	var cars []models.Car
	DB.Find(&cars)
	return cars
}

func CreateCar(car models.Car) {
	DB.Create(&car)
}

func UpdateCar(id string, car models.Car) {
	DB.Model(&models.Car{}).Where("id = ?", id).Updates(car)
}

func DeleteCar(id string) {
	DB.Delete(&models.Car{}, id)
}

func GetCarByID(id string) (*models.Car, error) {
	var car models.Car
	result := DB.First(&car, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &car, nil
}
