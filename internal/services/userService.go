package services

import (
	"Cars/internal/models"
)

func GetUsers() []models.User {
	var users []models.User
	DB.Find(&users)
	return users
}

func CreateUser(user models.User) {
	DB.Create(&user)
}

func UpdateUser(id string, user models.User) {
	DB.Model(&models.User{}).Where("id = ?", id).Updates(user)
}

func DeleteUser(id string) {
	DB.Delete(&models.User{}, id)
}
