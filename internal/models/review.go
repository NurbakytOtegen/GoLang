package models

import (
	"errors"
	"time"
)

// Review represents a car review with rating and comment
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CarID     uint      `json:"car_id" gorm:"index"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Rating    float32   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Car       Car       `json:"car" gorm:"foreignKey:CarID"`
}

// Validate performs basic validation of review data
func (r *Review) Validate() error {
	if r.CarID == 0 {
		return ErrCarIDRequired
	}
	if r.UserID == 0 {
		return ErrUserIDRequired
	}
	if r.Rating < 1 || r.Rating > 5 {
		return ErrInvalidRating
	}
	if len(r.Comment) < 3 {
		return ErrCommentTooShort
	}
	if len(r.Comment) > 1000 {
		return ErrCommentTooLong
	}
	return nil
}

// Custom errors for review validation
var (
	ErrCarIDRequired   = errors.New("car ID is required")
	ErrUserIDRequired  = errors.New("user ID is required")
	ErrInvalidRating   = errors.New("rating must be between 1 and 5")
	ErrCommentTooShort = errors.New("comment must be at least 3 characters long")
	ErrCommentTooLong  = errors.New("comment must not exceed 1000 characters")
)
