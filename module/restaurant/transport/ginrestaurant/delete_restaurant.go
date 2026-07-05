package ginrestaurant

import (
	restaurantbiz "food-delivery/module/restaurant/biz"
	restaurantstorage "food-delivery/module/restaurant/storage"
	"net/http" // HTTP status code constants (200, 400, ...)
	"strconv"

	"github.com/gin-gonic/gin" // Gin web framework
	"gorm.io/gorm"
)

func DeleteRestaurant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))


		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			})
			return
		}

		store := restaurantstorage.NewSQLStore(db)
		biz := restaurantbiz.NewDeleteRestaurantBiz(store)

		if err := biz.DeleteRestaurant(c.Request.Context(), id); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": 1,
		})
	}
}