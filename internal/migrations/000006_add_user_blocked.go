package migrations

import (
	"Cars/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// AddUserBlocked adds is_blocked field to User table
func AddUserBlocked() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "000006_add_user_blocked",
		Migrate: func(tx *gorm.DB) error {
			// Проверяем, существует ли колонка is_blocked
			if !tx.Migrator().HasColumn(&models.User{}, "is_blocked") {
				// Добавляем колонку is_blocked, если её нет
				if err := tx.Migrator().AddColumn(&models.User{}, "is_blocked"); err != nil {
					return err
				}

				// Устанавливаем значение по умолчанию false для существующих записей
				if err := tx.Model(&models.User{}).Where("is_blocked IS NULL").Update("is_blocked", false).Error; err != nil {
					return err
				}
			}

			// Проверяем наличие SUPER_ADMIN
			var count int64
			if err := tx.Model(&models.User{}).Where("role = ?", string(models.RoleSuperAdmin)).Count(&count).Error; err != nil {
				return err
			}

			// Создаем SUPER_ADMIN если его нет
			if count == 0 {
				defaultAdmin := &models.User{
					Name:      "Super Admin",
					Email:     "admin@admin.com",
					Password:  "$2a$10$ZGX4.UZ6NUxXfGt.4Zbh9.8NX.FzlHf8/dR.9A.UQ8YX3NqVn3IGm", // password: admin123
					Role:      string(models.RoleSuperAdmin),
					IsBlocked: false,
				}
				if err := tx.Create(defaultAdmin).Error; err != nil {
					return err
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Удаляем колонку is_blocked при откате
			if tx.Migrator().HasColumn(&models.User{}, "is_blocked") {
				if err := tx.Migrator().DropColumn(&models.User{}, "is_blocked"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
