package routes

import (
	"Cars/internal/controllers"
	"Cars/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupFavoriteRoutes(router *gin.Engine, favoriteController *controllers.FavoriteController) {
	// Группа маршрутов для работы с избранными автомобилями
	favorites := router.Group("/api/favorites")
	favorites.Use(middleware.AuthMiddleware()) // Защищаем все маршруты авторизацией

	{
		// Получить список избранных автомобилей пользователя
		favorites.GET("", favoriteController.GetUserFavorites)

		// Добавить автомобиль в избранное
		favorites.POST("/:car_id", favoriteController.AddToFavorites)

		// Удалить автомобиль из избранного
		favorites.DELETE("/:car_id", favoriteController.RemoveFromFavorites)

		// Проверить, находится ли автомобиль в избранном
		favorites.GET("/:car_id", favoriteController.IsFavorite)
	}
}
