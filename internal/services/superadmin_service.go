package services

import (
	"Cars/internal/models"
	"errors"

	"gorm.io/gorm"
)

type SuperAdminService struct {
	db *gorm.DB
}

func NewSuperAdminService(db *gorm.DB) *SuperAdminService {
	return &SuperAdminService{db: db}
}

// GetAllUsers returns all users except super admins
func (s *SuperAdminService) GetAllUsers() ([]models.UserResponse, error) {
	var users []models.User
	if err := s.db.Where("role != ?", string(models.RoleSuperAdmin)).Find(&users).Error; err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToUserResponse()
	}
	return userResponses, nil
}

// UpdateUserRole updates a user's role if the target is not a super admin
func (s *SuperAdminService) UpdateUserRole(userID uint, newRole string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	if user.Role == string(models.RoleSuperAdmin) {
		return errors.New("cannot modify super admin role")
	}

	// Validate new role
	if newRole != string(models.RoleUser) && newRole != string(models.RoleAdmin) {
		return errors.New("invalid role")
	}

	return s.db.Model(&user).Update("role", newRole).Error
}

// BlockUser blocks a user if they are not a super admin
func (s *SuperAdminService) BlockUser(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	if user.Role == string(models.RoleSuperAdmin) {
		return errors.New("cannot block super admin")
	}

	return s.db.Model(&user).Update("is_blocked", true).Error
}

// UnblockUser unblocks a user
func (s *SuperAdminService) UnblockUser(userID uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("is_blocked", false).Error
}

// GetUserByID gets a user by ID
func (s *SuperAdminService) GetUserByID(userID uint) (*models.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	response := user.ToUserResponse()
	return &response, nil
}
