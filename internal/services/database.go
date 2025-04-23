package services

import (
	"Cars/internal/migrations"
	_ "Cars/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=21012023 dbname=Carss port=5432 sslmode=disable TimeZone=Asia/Almaty"
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
