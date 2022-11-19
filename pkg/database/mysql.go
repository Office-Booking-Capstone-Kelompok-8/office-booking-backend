package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(dbHost string, dbPort string, dbUsername string, dbPassword string, dbName string) (*gorm.DB, error) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database Connected!")

	return DB, err
}
