package controllers

import (
	"net/http"
	"strconv"

	"Cars/internal/services"

	"github.com/gin-gonic/gin"
)

type FavoriteController struct {
	favoriteService *services.FavoriteService
}

func NewFavoriteController(favoriteService *services.FavoriteService) *FavoriteController {
	return &FavoriteController{
		favoriteService: favoriteService,
	}
}

// AddToFavorites обработчик для добавления автомобиля в избранное
func (c *FavoriteController) AddToFavorites(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	carID, err := strconv.ParseUint(ctx.Param("car_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	favorite, err := c.favoriteService.AddToFavorites(userID, uint(carID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, favorite)
}

// RemoveFromFavorites обработчик для удаления автомобиля из избранного
func (c *FavoriteController) RemoveFromFavorites(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	carID, err := strconv.ParseUint(ctx.Param("car_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	if err := c.favoriteService.RemoveFromFavorites(userID, uint(carID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "removed from favorites"})
}

// GetUserFavorites обработчик для получения списка избранных автомобилей
func (c *FavoriteController) GetUserFavorites(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	favorites, err := c.favoriteService.GetUserFavorites(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, favorites)
}

// IsFavorite обработчик для проверки, находится ли автомобиль в избранном
func (c *FavoriteController) IsFavorite(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	carID, err := strconv.ParseUint(ctx.Param("car_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	isFavorite := c.favoriteService.IsFavorite(userID, uint(carID))
	ctx.JSON(http.StatusOK, gin.H{"is_favorite": isFavorite})
}
