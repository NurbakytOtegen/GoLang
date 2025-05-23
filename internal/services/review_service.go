package services

import (
	"Cars/internal/models"
	"errors"

	"gorm.io/gorm"
)

// ReviewService handles business logic related to car reviews
type ReviewService struct {
	db *gorm.DB
}

// NewReviewService creates a new instance of ReviewService
func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

// Common errors
var (
	ErrReviewNotFound = errors.New("review not found")
	ErrInvalidReview  = errors.New("invalid review data")
)

// CreateReview creates a new review and updates car's average rating
func (s *ReviewService) CreateReview(review *models.Review) error {
	if err := review.Validate(); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create the review
		if err := tx.Create(review).Error; err != nil {
			return err
		}

		// Update car's average rating
		var avgRating float32
		if err := tx.Model(&models.Review{}).
			Where("car_id = ?", review.CarID).
			Select("COALESCE(AVG(rating), 0)").
			Scan(&avgRating).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Car{}).
			Where("id = ?", review.CarID).
			Update("avg_rating", avgRating).Error; err != nil {
			return err
		}

		// Load the car data for the response
		if err := tx.Preload("Car").First(review, review.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetReviewByID retrieves a review by its ID
func (s *ReviewService) GetReviewByID(id uint) (*models.Review, error) {
	var review models.Review
	if err := s.db.Preload("Car").First(&review, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return &review, nil
}

// UpdateReview updates an existing review and recalculates car's average rating
func (s *ReviewService) UpdateReview(review *models.Review) error {
	if err := review.Validate(); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update the review
		if err := tx.Save(review).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrReviewNotFound
			}
			return err
		}

		// Recalculate car's average rating
		var avgRating float32
		if err := tx.Model(&models.Review{}).
			Where("car_id = ?", review.CarID).
			Select("COALESCE(AVG(rating), 0)").
			Scan(&avgRating).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Car{}).
			Where("id = ?", review.CarID).
			Update("avg_rating", avgRating).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteReview removes a review and updates car's average rating
func (s *ReviewService) DeleteReview(id uint) error {
	var review models.Review
	if err := s.db.First(&review, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete the review
		if err := tx.Delete(&review).Error; err != nil {
			return err
		}

		// Recalculate car's average rating
		var avgRating float32
		if err := tx.Model(&models.Review{}).
			Where("car_id = ?", review.CarID).
			Select("COALESCE(AVG(rating), 0)").
			Scan(&avgRating).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Car{}).
			Where("id = ?", review.CarID).
			Update("avg_rating", avgRating).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetReviewsByCar retrieves all reviews for a specific car
func (s *ReviewService) GetReviewsByCar(carID uint) ([]models.Review, error) {
	var reviews []models.Review
	if err := s.db.Preload("Car").Where("car_id = ?", carID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// GetReviewsByUser retrieves all reviews by a specific user
func (s *ReviewService) GetReviewsByUser(userID uint) ([]models.Review, error) {
	var reviews []models.Review
	if err := s.db.Where("user_id = ?", userID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// GetCarRatingStats returns detailed rating statistics for a car
type RatingStats struct {
	AvgRating    float32 `json:"avg_rating"`
	TotalReviews int64   `json:"total_reviews"`
	Rating5Count int64   `json:"rating_5_count"`
	Rating4Count int64   `json:"rating_4_count"`
	Rating3Count int64   `json:"rating_3_count"`
	Rating2Count int64   `json:"rating_2_count"`
	Rating1Count int64   `json:"rating_1_count"`
}

func (s *ReviewService) GetCarRatingStats(carID uint) (*RatingStats, error) {
	var stats RatingStats

	// Get average rating and total reviews
	err := s.db.Model(&models.Review{}).
		Where("car_id = ?", carID).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as total_reviews").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// Get count for each rating
	for rating := 1; rating <= 5; rating++ {
		var count int64
		err := s.db.Model(&models.Review{}).
			Where("car_id = ? AND rating = ?", carID, rating).
			Count(&count).Error
		if err != nil {
			return nil, err
		}

		switch rating {
		case 1:
			stats.Rating1Count = count
		case 2:
			stats.Rating2Count = count
		case 3:
			stats.Rating3Count = count
		case 4:
			stats.Rating4Count = count
		case 5:
			stats.Rating5Count = count
		}
	}

	return &stats, nil
}

// GetTopRatedCars returns a list of cars sorted by average rating
func (s *ReviewService) GetTopRatedCars(limit int) ([]models.Car, error) {
	var cars []models.Car
	if err := s.db.Order("avg_rating DESC, id ASC").
		Where("brand != '' AND model != ''"). // Добавляем фильтр для непустых значений
		Limit(limit).
		Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}
