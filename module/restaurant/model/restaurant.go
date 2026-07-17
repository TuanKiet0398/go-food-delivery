package restaurantmodel

import "food-delivery/common"

const EntityName = "restaurants"

// Restaurant represents a row in the "restaurants" table
type Restaurant struct {
	common.SQLModel `json:",inline"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

// Restaurant represents a row in the "restaurants" table
type RestaurantCreate struct {
	common.SQLModel `json:",inline"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

// RestaurantUpdate is used for partial updates (PATCH); fields are pointers
// so we can distinguish "not provided" (nil) from "provided as empty value"
type RestaurantUpdate struct {
	common.SQLModel `json:",inline"`
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"addr" gorm:"column:addr;"`
}

// TableName tells GORM which table to use for Restaurant
func (Restaurant) TableName() string {
	return "restaurants"
}

func (r *Restaurant) Mask(isAdminOrOwner bool) {
	r.GenUID(common.DbTypeRestaurant)
}

func (data *RestaurantCreate) Mask(isAdminOrOwner bool) {
	data.GenUID(common.DbTypeRestaurant)
}

// TableName makes RestaurantUpdate share the same table as Restaurant
func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}


// TableName makes RestaurantUpdate share the same table as Restaurant
func (RestaurantCreate) TableName() string {
	return Restaurant{}.TableName()
}