package models

import "time"

// Favorite представляет избранный автомобиль пользователя
type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CarID     uint      `json:"car_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Car       Car       `json:"car" gorm:"foreignKey:CarID"`
}

// FavoriteResponse представляет ответ API для избранного автомобиля
type FavoriteResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Car       Car       `json:"car"`
	CreatedAt time.Time `json:"created_at"`
}
