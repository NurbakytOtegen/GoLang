package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// FixReviewsCascade fixes the cascade delete for reviews
func FixReviewsCascade() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "000007_fix_reviews_cascade",
		Migrate: func(tx *gorm.DB) error {
			// Сначала удаляем существующее ограничение
			if err := tx.Exec(`ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_cars_reviews;`).Error; err != nil {
				return err
			}

			// Создаем новое ограничение с CASCADE
			if err := tx.Exec(`
				ALTER TABLE reviews 
				ADD CONSTRAINT fk_cars_reviews 
				FOREIGN KEY (car_id) 
				REFERENCES cars(id) 
				ON DELETE CASCADE;
			`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// При откате создаем ограничение без CASCADE
			if err := tx.Exec(`ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_cars_reviews;`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				ALTER TABLE reviews 
				ADD CONSTRAINT fk_cars_reviews 
				FOREIGN KEY (car_id) 
				REFERENCES cars(id);
			`).Error; err != nil {
				return err
			}

			return nil
		},
	}
}
