package migrations

import (
	"Cars/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// AddCarType adds car_type field to cars table
func AddCarType() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "000004_add_car_type",
		Migrate: func(tx *gorm.DB) error {
			// Add car_type to cars table if it doesn't exist
			if !tx.Migrator().HasColumn(&models.Car{}, "car_type") {
				if err := tx.Migrator().AddColumn(&models.Car{}, "car_type"); err != nil {
					return err
				}
				// Set default value for existing records
				if err := tx.Exec(`UPDATE cars SET car_type = 'sedan' WHERE car_type IS NULL`).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove car_type from cars table if it exists
			if tx.Migrator().HasColumn(&models.Car{}, "car_type") {
				if err := tx.Migrator().DropColumn(&models.Car{}, "car_type"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
