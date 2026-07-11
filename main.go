package main

import (
	"food-delivery/component/appctx"
	"food-delivery/module/restaurant/transport/ginrestaurant"
	"log"      // logging to console
	"net/http" // HTTP status code constants (200, 400, ...)
	"os"       // read environment variables
	"strconv"  // convert string to number

	"github.com/gin-gonic/gin" // Gin web framework
	"github.com/joho/godotenv" // load .env file into environment variables
	"gorm.io/driver/mysql"     // MySQL driver for GORM
	"gorm.io/gorm"             // ORM for database access
)

// Restaurant represents a row in the "restaurants" table
type Restaurant struct {
	Id   int    `js:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
	Logo string `json:"logo" gorm:"column:logo;"`
}

// RestaurantUpdate is used for partial updates (PATCH); fields are pointers
// so we can distinguish "not provided" (nil) from "provided as empty value"
type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"addr" gorm:"column:addr;"`
}

// TableName tells GORM which table to use for Restaurant
func (Restaurant) TableName() string {
	return "restaurants"
}

// TableName makes RestaurantUpdate share the same table as Restaurant
func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	// Read the MySQL connection string and open a GORM connection
	dsn := os.Getenv("MYSQL_CONN_STRING")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	// Create a Gin router with default middleware (logger + recovery)
	r := gin.Default()

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	appContext := appctx.NewAppContext(db)

	// Route group for /v1/restaurants
	v1 := r.Group("/v1")
	restaurants := v1.Group("/restaurants")

	// Create a new restaurant
	restaurants.POST("", ginrestaurant.CreateRestaurant(appContext))

	// Get a single restaurant by id
	restaurants.GET("/:id", func(c *gin.Context) {
		// Parse the id path param into an int
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		var data Restaurant
		// Fetch the first record matching the given id
		db.Where("id = ?", id).First(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	// List restaurants with pagination
	restaurants.GET("", ginrestaurant.ListRestaurant(appContext))

	// Partially update a restaurant by id
	restaurants.PATCH("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		var data RestaurantUpdate

		// Only fields present in the request body get bound (others stay nil)
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Update the non-nil fields on the record matching the given id
		db.Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	// Delete a restaurant by id
	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))

	// Start the server, listening on port 8080 by default
	r.Run()

}
