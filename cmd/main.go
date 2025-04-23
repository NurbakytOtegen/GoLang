package main

import (
	"Cars/internal/controllers"
	"Cars/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// Мәліметтер базасына қосылу
	services.ConnectDatabase()

	// Gin роутерін құру
	router := gin.Default()

	// Маршруттарды тіркеу. Құрылғыны (router) функцияға өткізіңіз.
	controllers.RegisterAuthRoutes(router)
	controllers.RegisterCarRoutes(router)
	controllers.RegisterUserRoutes(router)

	//серверді іске қосамыз
	server := &http.Server{
		Addr:    ":8081", // мысалы, 8081
		Handler: router,
	}

	log.Println("Server running on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
