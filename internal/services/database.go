package services

import (
	_ "Cars/internal/models"
	"log"

	"Cars/internal/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=Cars port=5432"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ ПостгреSQL базасына қосыла алмадық: ", err)
	}

	// Gormigrate пайдалану
	m := migrations.GetMigrations(DB)

	if err := m.Migrate(); err != nil {
		log.Fatal("❌ Миграция қатесі: ", err)
	}

	log.Println("✅ PostgreSQL базасына сәтті қосылдық және миграциялар орындалды.")
}
