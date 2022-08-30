package services

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func ConnectDB() {

	d, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_CONNECTION")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting with database : %v", err)
	}

	db = d

}

func GetDB() *gorm.DB {
	return db
}
