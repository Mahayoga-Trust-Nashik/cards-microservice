package services

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func ConnectDB() {

	d, err := gorm.Open(mysql.Open("mahayoga_api:Ndk@1940knk@tcp(mahayoga-database:3306)/mahayoga_mobile?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting with database : %v", err)
	}

	db = d

}

func GetDB() *gorm.DB {
	return db
}
