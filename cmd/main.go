package main

import (
	"Cars/internal/controllers"
	"Cars/internal/migrations"
	"Cars/internal/routes"
	"Cars/internal/services"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Мәліметтер базасына қосылу
	services.ConnectDatabase()

	//migrate
	migrate := migrations.GetMigrations(services.DB)
	if err := migrate.Migrate(); err != nil {
		log.Fatal("Migration failed: %v", err)
	}
	log.Println("Migrations applied successfully")

	// Gin роутерін құру
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize services
	reviewService := services.NewReviewService(services.DB)
	profileService := services.NewProfileService(services.DB)
	superAdminService := services.NewSuperAdminService(services.DB)
	favoriteService := services.NewFavoriteService(services.DB)

	// Initialize controllers
	reviewController := controllers.NewReviewController(reviewService)
	profileController := controllers.NewProfileController(profileService)
	superAdminController := controllers.NewSuperAdminController(superAdminService)
	favoriteController := controllers.NewFavoriteController(favoriteService)

	// Маршруттарды тіркеу
	controllers.RegisterAuthRoutes(router)
	controllers.RegisterCarRoutes(router)
	controllers.RegisterUserRoutes(router)
	routes.SetupReviewRoutes(router, reviewController)
	controllers.RegisterProfileRoutes(router, profileController)
	controllers.RegisterSuperAdminRoutes(router, superAdminController)
	routes.SetupFavoriteRoutes(router, favoriteController)

	//серверді іске қосамыз
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Println("Server running on http://localhost:8081")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
