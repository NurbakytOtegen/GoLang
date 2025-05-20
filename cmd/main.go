package main

import (
	"Cars/internal/controllers"
	"Cars/internal/migrations"
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
	log.Println("Migartions applied successfully")

	// Gin роутерін құру
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Маршруттарды тіркеу. Құрылғыны (router) функцияға өткізіңіз.
	controllers.RegisterAuthRoutes(router)
	controllers.RegisterCarRoutes(router)
	controllers.RegisterUserRoutes(router)

	//серверді іске қосамыз
	server := &http.Server{
		Addr:    ":8081", // мысалы, 8081
		Handler: router,
	}

	log.Println("Server running on http://localhost:8081")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
