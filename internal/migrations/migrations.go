package migrations

import (
	"Cars/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// GetMigrations returns all migrations in the correct order
func GetMigrations(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20250412_create_users_and_cars",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{}, &models.Car{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users", "cars")
			},
		},
		AddReviewsAndRatings(),
		AddCarTypeGorm(),
		AddUserBlocked(),
		FixReviewsCascade(),
		AddFavorites(),
	})
}
