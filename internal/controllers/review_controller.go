package controllers

import (
	"Cars/internal/models"
	"Cars/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	reviewService *services.ReviewService
}

func NewReviewController(reviewService *services.ReviewService) *ReviewController {
	return &ReviewController{
		reviewService: reviewService,
	}
}

// CreateReview godoc
// @Summary Create a new review
// @Description Create a new review for a car
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.Review true "Review object"
// @Success 201 {object} models.Review
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reviews [post]
func (c *ReviewController) CreateReview(ctx *gin.Context) {
	var review models.Review
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Get user ID from context (assuming it was set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}
	review.UserID = userID.(uint)

	if err := c.reviewService.CreateReview(&review); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, review)
}

// GetReview godoc
// @Summary Get a review by ID
// @Description Get a review's details by its ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} models.Review
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reviews/{id} [get]
func (c *ReviewController) GetReview(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid review ID"})
		return
	}

	review, err := c.reviewService.GetReviewByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrReviewNotFound {
			status = http.StatusNotFound
		}
		ctx.JSON(status, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, review)
}

// UpdateReview godoc
// @Summary Update a review
// @Description Update an existing review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body models.Review true "Review object"
// @Success 200 {object} models.Review
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reviews/{id} [put]
func (c *ReviewController) UpdateReview(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid review ID"})
		return
	}

	var review models.Review
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Verify user owns this review
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	existingReview, err := c.reviewService.GetReviewByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrReviewNotFound {
			status = http.StatusNotFound
		}
		ctx.JSON(status, ErrorResponse{Error: err.Error()})
		return
	}

	if existingReview.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, ErrorResponse{Error: "Not authorized to update this review"})
		return
	}

	review.ID = uint(id)
	review.UserID = userID.(uint)

	if err := c.reviewService.UpdateReview(&review); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, review)
}

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete a review by ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reviews/{id} [delete]
func (c *ReviewController) DeleteReview(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid review ID"})
		return
	}

	// Verify user owns this review
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	existingReview, err := c.reviewService.GetReviewByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrReviewNotFound {
			status = http.StatusNotFound
		}
		ctx.JSON(status, ErrorResponse{Error: err.Error()})
		return
	}

	if existingReview.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, ErrorResponse{Error: "Not authorized to delete this review"})
		return
	}

	if err := c.reviewService.DeleteReview(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetCarReviews godoc
// @Summary Get reviews for a car
// @Description Get all reviews for a specific car
// @Tags reviews
// @Produce json
// @Param id path int true "Car ID"
// @Success 200 {array} models.Review
// @Failure 500 {object} ErrorResponse
// @Router /cars/{id}/reviews [get]
func (c *ReviewController) GetCarReviews(ctx *gin.Context) {
	carID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid car ID"})
		return
	}

	reviews, err := c.reviewService.GetReviewsByCar(uint(carID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

// GetCarStats godoc
// @Summary Get car rating statistics
// @Description Get detailed rating statistics for a car
// @Tags reviews
// @Produce json
// @Param id path int true "Car ID"
// @Success 200 {object} services.RatingStats
// @Failure 500 {object} ErrorResponse
// @Router /cars/{id}/rating-stats [get]
func (c *ReviewController) GetCarStats(ctx *gin.Context) {
	carID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid car ID"})
		return
	}

	stats, err := c.reviewService.GetCarRatingStats(uint(carID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetTopRatedCars godoc
// @Summary Get top rated cars
// @Description Get a list of cars sorted by average rating
// @Tags reviews
// @Produce json
// @Param limit query int false "Number of cars to return" default(10)
// @Success 200 {array} models.Car
// @Failure 500 {object} ErrorResponse
// @Router /reviews/top-rated-cars [get]
func (c *ReviewController) GetTopRatedCars(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	cars, err := c.reviewService.GetTopRatedCars(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cars)
}
