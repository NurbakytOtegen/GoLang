package services

import (
	"Cars/internal/models"
	"errors"

	"gorm.io/gorm"
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

func DeleteCar(id string) error {
	// Проверяем существование машины
	var car models.Car
	if err := DB.First(&car, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("car not found")
		}
		return err
	}

	// Удаляем машину
	result := DB.Delete(&car)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("failed to delete car")
	}

	return nil
}

func GetCarByID(id string) (*models.Car, error) {
	var car models.Car
	result := DB.First(&car, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &car, nil
}
