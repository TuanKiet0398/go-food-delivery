package restaurantmodel


// Restaurant represents a row in the "restaurants" table
type Restaurant struct {
	Id   int    `js:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
	Status string `json:"status" gorm:"column:status;"`
}

// Restaurant represents a row in the "restaurants" table
type RestaurantCreate struct {
	Id   int    `js:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
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


// TableName makes RestaurantUpdate share the same table as Restaurant
func (RestaurantCreate) TableName() string {
	return Restaurant{}.TableName()
}