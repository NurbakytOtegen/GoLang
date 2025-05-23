package services

import (
	"Cars/internal/models"
	"errors"

	"gorm.io/gorm"
)

type FavoriteService struct {
	db *gorm.DB
}

func NewFavoriteService(db *gorm.DB) *FavoriteService {
	return &FavoriteService{db: db}
}

// AddToFavorites добавляет автомобиль в избранное пользователя
func (s *FavoriteService) AddToFavorites(userID, carID uint) (*models.Favorite, error) {
	// Проверяем, существует ли уже такая запись
	var existing models.Favorite
	result := s.db.Where("user_id = ? AND car_id = ?", userID, carID).First(&existing)
	if result.Error == nil {
		return nil, errors.New("car is already in favorites")
	}

	favorite := &models.Favorite{
		UserID: userID,
		CarID:  carID,
	}

	if err := s.db.Create(favorite).Error; err != nil {
		return nil, err
	}

	// Загружаем связанные данные
	s.db.Preload("Car").First(favorite, favorite.ID)
	return favorite, nil
}

// RemoveFromFavorites удаляет автомобиль из избранного пользователя
func (s *FavoriteService) RemoveFromFavorites(userID, carID uint) error {
	result := s.db.Where("user_id = ? AND car_id = ?", userID, carID).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// GetUserFavorites получает список избранных автомобилей пользователя
func (s *FavoriteService) GetUserFavorites(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := s.db.Where("user_id = ?", userID).Preload("Car").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

// IsFavorite проверяет, находится ли автомобиль в избранном у пользователя
func (s *FavoriteService) IsFavorite(userID, carID uint) bool {
	var favorite models.Favorite
	result := s.db.Where("user_id = ? AND car_id = ?", userID, carID).First(&favorite)
	return result.Error == nil
}
