package main

import (
	"log"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/database/mysql"
	"office-booking-backend/pkg/entity"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("ENV") == "production" {
		return
	}

	//	load env variables from .env file for local development
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	env := config.LoadConfig()

	db := mysql.InitDatabase(env["DB_HOST"], env["DB_PORT"], env["DB_USER"], env["DB_PASS"], env["DB_NAME"])

	err := db.AutoMigrate(
		&entity.User{},
		&entity.UserDetail{},
		&entity.ProfilePicture{},
	)

	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	log.Println("Database migration successful")
}
