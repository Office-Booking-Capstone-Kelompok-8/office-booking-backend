package main

import (
	"log"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/database/mysql"
	"office-booking-backend/pkg/entity"
	"office-booking-backend/pkg/utils/password"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
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

	err = InitCity(db)
	if err != nil {
		log.Fatalf("Error seeding city: %v", err)
	}

	err = InitDistrict(db)
	if err != nil {
		log.Fatalf("Error seeding district: %v", err)
	}

	err = InitCategory(db)
	if err != nil {
		log.Fatalf("Error seeding category: %v", err)
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
			Message: "Active",
		},
		{
			ID:      7,
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

func InitCity(db *gorm.DB) error {
	//  ID, Name
	// (153,'Kab. Kepulauan Seribu'),
	// (154,'Kota Jakarta Selatan'),
	// (155,'Kota Jakarta Timur'),
	// (156,'Kota Jakarta Pusat'),
	// (157,'Kota Jakarta Barat'),
	// (158,'Kota Jakarta Utara'),

	city := entity.Cities{
		{
			ID:   153,
			Name: "Kab. Kepulauan Seribu",
		},
		{
			ID:   154,
			Name: "Kota Jakarta Selatan",
		},
		{
			ID:   155,
			Name: "Kota Jakarta Timur",
		},
		{
			ID:   156,
			Name: "Kota Jakarta Pusat",
		},
		{
			ID:   157,
			Name: "Kota Jakarta Barat",
		},
		{
			ID:   158,
			Name: "Kota Jakarta Utara",
		},
	}

	// check if city already exists
	var count int64
	db.Model(&entity.City{}).Count(&count)
	if count > 0 {
		return nil
	}

	return db.Create(&city).Error
}

func InitDistrict(db *gorm.DB) error {
	//  ID, CityID, Name
	// (1880,153,' Kepulauan Seribu Selatan'),
	// (1881,153,' Kepulauan Seribu Utara'),
	// (1882,154,' Jagakarsa'),
	// (1883,154,' Pasar Minggu'),
	// (1884,154,' Cilandak'),
	// (1885,154,' Pesanggrahan'),
	// (1886,154,' Kebayoran Lama'),
	// (1887,154,' Kebayoran Baru'),
	// (1888,154,' Mampang Prapatan'),
	// (1889,154,' Pancoran'),
	// (1890,154,' Tebet'),
	// (1891,154,' Setia Budi'),
	// (1892,155,' Pasar Rebo'),
	// (1893,155,' Ciracas'),
	// (1894,155,' Cipayung'),
	// (1895,155,' Makasar'),
	// (1896,155,' Kramat Jati'),
	// (1897,155,' Jatinegara'),
	// (1898,155,' Duren Sawit'),
	// (1899,155,' Cakung'),
	// (1900,155,' Pulo Gadung'),
	// (1901,155,' Matraman'),
	// (1902,156,' Tanah Abang'),
	// (1903,156,' Menteng'),
	// (1904,156,' Senen'),
	// (1905,156,' Johar Baru'),
	// (1906,156,' Cempaka Putih'),
	// (1907,156,' Kemayoran'),
	// (1908,156,' Sawah Besar'),
	// (1909,156,' Gambir'),
	// (1910,157,' Kembangan'),
	// (1911,157,' Kebon Jeruk'),
	// (1912,157,' Palmerah'),
	// (1913,157,' Grogol Petamburan'),
	// (1914,157,' Tambora'),
	// (1915,157,' Taman Sari'),
	// (1916,157,' Cengkareng'),
	// (1917,157,' Kali Deres'),
	// (1918,158,' Penjaringan'),
	// (1919,158,' Pademangan'),
	// (1920,158,' Tanjung Priok'),
	// (1921,158,' Koja'),
	// (1922,158,' Kelapa Gading'),
	// (1923,158,' Cilincing'),

	district := entity.Districts{
		{
			ID:     1880,
			CityID: 153,
			Name:   "Kepulauan Seribu Selatan",
		},
		{
			ID:     1881,
			CityID: 153,
			Name:   "Kepulauan Seribu Utara",
		},
		{
			ID:     1882,
			CityID: 154,
			Name:   "Jagakarsa",
		},
		{
			ID:     1883,
			CityID: 154,
			Name:   "Pasar Minggu",
		},
		{
			ID:     1884,
			CityID: 154,
			Name:   "Cilandak",
		},
		{
			ID:     1885,
			CityID: 154,
			Name:   "Pesanggrahan",
		},
		{
			ID:     1886,
			CityID: 154,
			Name:   "Kebayoran Lama",
		},
		{
			ID:     1887,
			CityID: 154,
			Name:   "Kebayoran Baru",
		},
		{
			ID:     1888,
			CityID: 154,
			Name:   "Mampang Prapatan",
		},
		{
			ID:     1889,
			CityID: 154,
			Name:   "Pancoran",
		},
		{
			ID:     1890,
			CityID: 154,
			Name:   "Tebet",
		},
		{
			ID:     1891,
			CityID: 154,
			Name:   "Setia Budi",
		},
		{
			ID:     1892,
			CityID: 155,
			Name:   "Pasar Rebo",
		},
		{
			ID:     1893,
			CityID: 155,
			Name:   "Ciracas",
		},
		{
			ID:     1894,
			CityID: 155,
			Name:   "Cipayung",
		},
		{

			ID:     1895,
			CityID: 155,
			Name:   "Makasar",
		},
		{
			ID:     1896,
			CityID: 155,
			Name:   "Kramat Jati",
		},
		{
			ID:     1897,
			CityID: 155,
			Name:   "Jatinegara",
		},
		{
			ID:     1898,
			CityID: 155,

			Name: "Duren Sawit",
		},
		{
			ID:     1899,
			CityID: 155,
			Name:   "Cakung",
		},
		{
			ID:     1900,
			CityID: 155,
			Name:   "Pulo Gadung",
		},
		{
			ID:     1901,
			CityID: 155,
			Name:   "Matraman",
		},
		{
			ID:     1902,
			CityID: 156,
			Name:   "Tanah Abang",
		},
		{
			ID:     1903,
			CityID: 156,
			Name:   "Menteng",
		},
		{
			ID:     1904,
			CityID: 156,
			Name:   "Senen",
		},
		{
			ID:     1905,
			CityID: 156,
			Name:   "Johar Baru",
		},
		{
			ID:     1906,
			CityID: 156,
			Name:   "Cempaka Putih",
		},
		{
			ID:     1907,
			CityID: 156,
			Name:   "Kemayoran",
		},
		{
			ID:     1908,
			CityID: 156,
			Name:   "Sawah Besar",
		},
		{
			ID:     1909,
			CityID: 156,
			Name:   "Gambir",
		},
		{
			ID:     1910,
			CityID: 157,
			Name:   "Kembangan",
		},
		{
			ID:     1911,
			CityID: 157,
			Name:   "Kebon Jeruk",
		},
		{
			ID:     1912,
			CityID: 157,
			Name:   "Palmerah",
		},
		{
			ID:     1913,
			CityID: 157,
			Name:   "Grogol Petamburan",
		},
		{
			ID:     1914,
			CityID: 157,
			Name:   "Tambora",
		},
		{
			ID:     1915,
			CityID: 157,
			Name:   "Taman Sari",
		},
		{
			ID:     1916,
			CityID: 157,
			Name:   "Cengkareng",
		},
		{
			ID:     1917,
			CityID: 157,
			Name:   "Kali Deres",
		},
		{
			ID:     1918,
			CityID: 158,
			Name:   "Penjaringan",
		},
		{
			ID:     1919,
			CityID: 158,
			Name:   "Pademangan",
		},
		{
			ID:     1920,
			CityID: 158,
			Name:   "Tanjung Priok",
		},
		{
			ID:     1921,
			CityID: 158,
			Name:   "Koja",
		},
		{
			ID:     1922,
			CityID: 158,
			Name:   "Kelapa Gading",
		},
		{
			ID:     1923,
			CityID: 158,
			Name:   "Cilincing",
		},
	}

	var count int64
	db.Model(&entity.District{}).Count(&count)
	if count != 0 {
		return nil
	}

	return db.Create(&district).Error
}

func InitCategory(db *gorm.DB) error {
	category := entity.Categories{
		{
			ID:   1,
			Name: "other_houses",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/other_houses_FILL0_wght400_GRAD0_opsz48_ySm9eMGlM.svg",
		},
	}

	var count int64
	db.Model(&entity.Category{}).Count(&count)
	if count != 0 {
		return nil
	}

	return db.Create(&category).Error
}
