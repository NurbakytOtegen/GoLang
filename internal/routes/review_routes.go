package routes

import (
	"Cars/internal/controllers"
	"Cars/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupReviewRoutes configures all routes for review operations
func SetupReviewRoutes(router *gin.Engine, reviewController *controllers.ReviewController) {
	// Public routes for reviews
	reviews := router.Group("/reviews")
	{
		reviews.GET("/:id", reviewController.GetReview)
		reviews.GET("/top-rated-cars", reviewController.GetTopRatedCars)
	}

	// Car-specific review routes
	cars := router.Group("/cars")
	{
		// Важно: эти маршруты должны быть определены до других маршрутов с :id
		cars.GET("/:id/reviews", reviewController.GetCarReviews)
		cars.GET("/:id/rating-stats", reviewController.GetCarStats)
	}

	// Protected routes (require authentication)
	authReviews := router.Group("/reviews").Use(middleware.AuthMiddleware())
	{
		authReviews.POST("", reviewController.CreateReview)
		authReviews.PUT("/:id", reviewController.UpdateReview)
		authReviews.DELETE("/:id", reviewController.DeleteReview)
	}
}
