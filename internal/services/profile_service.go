package services

import (
	"Cars/internal/models"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// ProfileService handles profile-related business logic
type ProfileService struct {
	db *gorm.DB
}

// NewProfileService creates a new instance of ProfileService
func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}

// GetUserProfile retrieves user profile information with their reviews
func (s *ProfileService) GetUserProfile(userID uint) (*models.UserResponse, []models.Review, error) {
	// Get user information
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, err
	}

	// Get user's reviews
	var reviews []models.Review
	err = s.db.Where("user_id = ?", userID).Preload("Car").Find(&reviews).Error
	if err != nil {
		return nil, nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, reviews, nil
}

// GetUserProfileByEmail retrieves user profile information with their reviews by email
func (s *ProfileService) GetUserProfileByEmail(email string) (*models.UserResponse, []models.Review, error) {
	// Get user information
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, err
	}

	// Get user's reviews
	var reviews []models.Review
	err = s.db.Where("user_id = ?", user.ID).Preload("Car").Find(&reviews).Error
	if err != nil {
		return nil, nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, reviews, nil
}

// GetUserProfileByName retrieves user profile information with their reviews by name
func (s *ProfileService) GetUserProfileByName(name string) (*models.UserResponse, []models.Review, error) {
	// Get user information
	var user models.User
	err := s.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, err
	}

	// Get user's reviews
	var reviews []models.Review
	err = s.db.Where("user_id = ?", user.ID).Preload("Car").Find(&reviews).Error
	if err != nil {
		return nil, nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, reviews, nil
}

// UpdateUserName updates the user's name
func (s *ProfileService) UpdateUserName(userID uint, newName string) (*models.UserResponse, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Name = newName
	err = s.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateUserPassword updates the user's password
func (s *ProfileService) UpdateUserPassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Проверяем старый пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("incorrect old password")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.db.Save(&user).Error
}
