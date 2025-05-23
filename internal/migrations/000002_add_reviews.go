package migrations

import (
	"Cars/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// AddReviewsAndRatings creates reviews table and adds rating fields to cars
func AddReviewsAndRatings() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "000002_add_reviews",
		Migrate: func(tx *gorm.DB) error {
			// Add avg_rating to cars table if it doesn't exist
			if !tx.Migrator().HasColumn(&models.Car{}, "avg_rating") {
				if err := tx.Migrator().AddColumn(&models.Car{}, "avg_rating"); err != nil {
					return err
				}
			}

			// Create reviews table if it doesn't exist
			if !tx.Migrator().HasTable(&models.Review{}) {
				if err := tx.AutoMigrate(&models.Review{}); err != nil {
					return err
				}

				// Add foreign key constraint
				if err := tx.Exec(`
					ALTER TABLE reviews 
					ADD CONSTRAINT fk_cars_reviews 
					FOREIGN KEY (car_id) 
					REFERENCES cars(id) 
					ON DELETE CASCADE
				`).Error; err != nil {
					return err
				}

				// Create indexes
				if err := tx.Exec(`
					CREATE INDEX IF NOT EXISTS idx_reviews_car_id ON reviews(car_id);
					CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
					CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
				`).Error; err != nil {
					return err
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Drop indexes if they exist
			if tx.Migrator().HasTable("reviews") {
				if err := tx.Exec(`
					DROP INDEX IF EXISTS idx_reviews_rating;
					DROP INDEX IF EXISTS idx_reviews_user_id;
					DROP INDEX IF EXISTS idx_reviews_car_id;
				`).Error; err != nil {
					return err
				}

				// Drop reviews table (this will also drop the foreign key constraint)
				if err := tx.Migrator().DropTable("reviews"); err != nil {
					return err
				}
			}

			// Remove avg_rating from cars table if it exists
			if tx.Migrator().HasColumn(&models.Car{}, "avg_rating") {
				if err := tx.Migrator().DropColumn(&models.Car{}, "avg_rating"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
