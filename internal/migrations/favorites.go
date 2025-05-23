package migrations

import (
	"Cars/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// AddFavorites creates a migration for the favorites table
func AddFavorites() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250524_add_favorites",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Favorite{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("favorites")
		},
	}
}
