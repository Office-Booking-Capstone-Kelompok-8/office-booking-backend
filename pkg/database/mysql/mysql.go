package mysql

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(dbHost string, dbPort string, dbUsername string, dbPassword string, dbName string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connected successfully")

	return db
}
