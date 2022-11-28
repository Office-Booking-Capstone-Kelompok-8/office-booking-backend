package main

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/database/mysql"
	"office-booking-backend/pkg/entity"
	"office-booking-backend/pkg/utils/password"
	"os"
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
		&entity.Category{},
		&entity.Facility{},
		&entity.City{},
		&entity.District{},
		&entity.Building{},
		&entity.Picture{},
		&entity.Payment{},
		&entity.PaymentPicture{},
		&entity.Status{},
		&entity.Reservation{},
	)

	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	err = InitAdmin(db)
	if err != nil {
		log.Fatalf("Error seeding admin: %v", err)
	}

	err = InitStatus(db)
	if err != nil {
		log.Fatalf("Error seeding status: %v", err)
	}

	log.Println("Database migration successful")
}

func InitAdmin(db *gorm.DB) error {
	passFunc := password.NewPasswordFuncImpl()
	pass, err := passFunc.GenerateFromPassword([]byte("admin123"), 10)
	if err != nil {
		return err
	}

	admin := entity.User{
		ID:         uuid.New().String(),
		Email:      "admin@mail.fortyfourvisual.com",
		Password:   string(pass), // admin123
		Role:       2,
		IsVerified: true,
		Detail: entity.UserDetail{
			Name:      "Admin",
			Phone:     "081234567890",
			PictureID: "123",
			Picture: entity.ProfilePicture{
				ID:  "123",
				Url: "https://ik.imagekit.io/fortyfour/default-image.jpg",
			},
		},
	}

	// check if admin already exists
	var count int64
	db.Model(&entity.User{}).Where("email = ?", admin.Email).Count(&count)
	if count > 0 {
		return nil
	}

	return db.Create(&admin).Error
}

func InitStatus(db *gorm.DB) error {
	status := []entity.Status{
		{
			ID:      1,
			Message: "Pending",
		},
		{
			ID:      2,
			Message: "Accepted",
		},
		{
			ID:      3,
			Message: "Rejected",
		},
		{
			ID:      4,
			Message: "Canceled",
		},
		{
			ID:      5,
			Message: "Awaiting Payment",
		},
		{
			ID:      6,
			Message: "Completed",
		},
	}

	// check if status already exists
	var count int64
	db.Model(&entity.Status{}).Count(&count)
	if count > 0 {
		return nil
	}

	return db.Create(&status).Error
}
