package ginrestaurant

import (
	"food-delivery/component/appctx"
	restaurantbiz "food-delivery/module/restaurant/biz"
	restaurantstorage "food-delivery/module/restaurant/storage"
	"net/http" // HTTP status code constants (200, 400, ...)
	"strconv"
	"food-delivery/common"
	"github.com/gin-gonic/gin" // Gin web framework

)

func DeleteRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		db := appCtx.GetMainDBConnection()

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

		c.JSON(200, common.SimpleSuccessResponse(true))

	}
}