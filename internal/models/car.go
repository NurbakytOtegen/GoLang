package models

import (
	"errors"
	"time"
)

type Car struct {
	ID           uint     `json:"id" gorm:"primaryKey"`
	Brand        string   `json:"brand"`
	Model        string   `json:"model"`
	CarType      string   `json:"car_type"`
	Year         int      `json:"year"`
	ImageURL     string   `json:"image_url"`
	Mileage      float64  `json:"mileage"`
	Transmission string   `json:"transmission"`
	EngineVolume float64  `json:"engine_vol"`
	Price        float64  `json:"price"`
	IsNew        bool     `json:"is_new"`
	AvgRating    float32  `json:"avg_rating" gorm:"default:0"`
	Reviews      []Review `json:"reviews,omitempty" gorm:"foreignKey:CarID"`
}

// CarStatus represents the current status of the car
type CarStatus string

const (
	StatusAvailable   CarStatus = "available"
	StatusSold        CarStatus = "sold"
	StatusMaintenance CarStatus = "maintenance"
	StatusReservation CarStatus = "reservation"
)

// Validate performs basic validation of car data
func (c *Car) Validate() error {
	if c.Brand == "" {
		return ErrEmptyBrand
	}
	if c.Model == "" {
		return ErrEmptyModel
	}
	if c.Year < 1900 || c.Year > time.Now().Year()+1 {
		return ErrInvalidYear
	}
	if c.EngineVolume <= 0 {
		return ErrInvalidEngineVolume
	}
	if c.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}

// Custom errors for car validation
var (
	ErrEmptyBrand          = errors.New("brand cannot be empty")
	ErrEmptyModel          = errors.New("model cannot be empty")
	ErrInvalidYear         = errors.New("invalid year")
	ErrInvalidEngineVolume = errors.New("invalid engine volume")
	ErrInvalidPrice        = errors.New("invalid price")
)
