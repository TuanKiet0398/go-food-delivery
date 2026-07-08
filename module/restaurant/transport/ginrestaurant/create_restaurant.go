package ginrestaurant

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	restaurantbiz "food-delivery/module/restaurant/biz"
	restaurantmodel "food-delivery/module/restaurant/model"
	restaurantstorage "food-delivery/module/restaurant/storage"
	"net/http" // HTTP status code constants (200, 400, ...)

	"github.com/gin-gonic/gin" // Gin web framework
)

func CreateRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		db := appCtx.GetMainDBConnection()

		var data restaurantmodel.RestaurantCreate

		// Bind JSON request body into data
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := restaurantstorage.NewSQLStore(db)
		biz := restaurantbiz.NewCreateRestaurantBiz(store)

		if err := biz.CreateRestaurant(c.Request.Context(), &data); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, common.SimpleSuccessResponse(data.Id))
	}
}