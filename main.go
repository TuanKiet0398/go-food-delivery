package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



type Restaurant struct {
	Id int `js:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
	Logo string `json:"logo" gorm:"column:logo;"`
}


func(Restaurant) TableName() string {
	return "restaurants"
}


func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	dsn := os.Getenv("MYSQL_CONN_STRING")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	newRestaurant := Restaurant{Name: "Kiet Ho", Addr: "267 Vuon Lai"}
	if err := db.Create(&newRestaurant).Error; err != nil {
		log.Println(err)
	}

	log.Println("New id", newRestaurant.Id)
}
