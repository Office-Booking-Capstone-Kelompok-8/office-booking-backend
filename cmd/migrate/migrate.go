package main

import (
	"flag"
	"fmt"
	"log"
	"office-booking-backend/pkg/constant"
	"office-booking-backend/pkg/custom"
	"office-booking-backend/pkg/database/mysql"
	"office-booking-backend/pkg/entity"
	"office-booking-backend/pkg/utils/password"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	var dbHost, dbPort, dbUser, dbPass, dbName, config string
	flag.StringVar(&dbHost, "host", "localhost", "Database host")
	flag.StringVar(&dbPort, "port", "3306", "Database port")
	flag.StringVar(&dbUser, "user", "root", "Database user")
	flag.StringVar(&dbPass, "pass", "root", "Database password")
	flag.StringVar(&dbName, "name", "", "Database name")
	flag.StringVar(&config, "config", "", "Config file path")
	flag.Parse()

	if config != "" {
		conf := viper.New()
		conf.SetConfigFile(config)

		err := conf.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		fmt.Println(conf.GetString("database.host"))

		dbHost = conf.GetString("service.db.host")
		dbPort = conf.GetString("service.db.port")
		dbUser = conf.GetString("service.db.user")
		dbPass = conf.GetString("service.db.pass")
		dbName = conf.GetString("service.db.name")
	}

	db := mysql.InitDatabase(
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	err := db.AutoMigrate(
		&entity.User{},
		&entity.UserDetail{},
		&entity.Building{},
		&entity.ProfilePicture{},
		&entity.Category{},
		&entity.Facility{},
		&entity.City{},
		&entity.District{},
		&entity.Picture{},
		&entity.Payment{},
		&entity.Bank{},
		&entity.Status{},
		&entity.Reservation{},
		&entity.Transaction{},
		&entity.Review{},
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

	err = InitBank(db)
	if err != nil {
		log.Fatalf("Error seeding bank: %v", err)
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
		IsVerified: custom.Bool(true),
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
			ID:      constant.PENDING_STATUS,
			Message: "Pending",
		},
		{
			ID:      constant.REJECTED_STATUS,
			Message: "Rejected",
		},
		{
			ID:      constant.CANCELED_STATUS,
			Message: "Canceled",
		},
		{
			ID:      constant.AWAITING_PAYMENT_STATUS,
			Message: "Awaiting Payment",
		},
		{
			ID:      constant.ACTIVE_STATUS,
			Message: "Active",
		},
		{
			ID:      constant.COMPLETED_STATUS,
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
			Name: "Lainnya",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/other_houses_FILL0_wght400_GRAD0_opsz48_ySm9eMGlM.svg",
		},
		{
			ID:   2,
			Name: "CCTV",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/device-cctv_BSxTxr0s1.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427963176",
		},
		{
			ID:   3,
			Name: "Dapur",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/kitchen_FILL0_wght400_GRAD0_opsz48_RzlEJkIJoC.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428432557",
		},
		{
			ID:   4,
			Name: "Kamar Mandi",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/Toilet_Bowl-595b40b85ba036ed117db567_Ie7dc45m5.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428666242",
		},
		{
			ID:   5,
			Name: "Parkir",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/local_parking_FILL0_wght400_GRAD0_opsz48_FgcoxojM2.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428432580",
		},
		{
			ID:   6,
			Name: "Wifi",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/wifi_FILL0_wght400_GRAD0_opsz48_OTrB_E3Ge.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428432560",
		},
		{
			ID:   7,
			Name: "Fitness",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/fitness_center_FILL0_wght400_GRAD0_opsz48_6vlA6Rv8m.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428433245",
		},
		{
			ID:   8,
			Name: "Telepon",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/call_FILL0_wght400_GRAD0_opsz48_yZP_MUvaNH.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428432624",
		},
		{
			ID:   9,
			Name: "Kantin",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/food_bank_FILL0_wght400_GRAD0_opsz48_hTRRszwP6.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428433169",
		},
		{
			ID:   10,
			Name: "AC",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/ac_unit_FILL0_wght400_GRAD0_opsz48_WgKuo7ULx.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428432656",
		},
		{
			ID:   11,
			Name: "Kipas Angin",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/mode_fan_FILL0_wght400_GRAD0_opsz48_UfKEzKxJ_T.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428432669",
		},
		{
			ID:   12,
			Name: "Televisi",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/desktop_mac_FILL0_wght400_GRAD0_opsz48_a-P5J2q_Eu.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428434101",
		},
		{
			ID:   13,
			Name: "Proyektor",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/projector-svgrepo-com_-I7UfrwD_.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428695681",
		},
		{
			ID:   14,
			Name: "Printer",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/printer-svgrepo-com_CwEwxlp9m.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670428760161",
		},
		{
			ID:   15,
			Name: "Ruang Merokok",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/smoking_rooms_FILL0_wght400_GRAD0_opsz48_gEsimHKIG.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670429102985",
		},
		{
			ID:   16,
			Name: "Petugas Kebersihan",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/mop_FILL0_wght400_GRAD0_opsz48_bJTBx-3VX.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670429103012",
		},
		{
			ID:   17,
			Name: "Meja",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/table_restaurant_FILL0_wght400_GRAD0_opsz48_1hyC5DKIp.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670429104525s",
		},
		{
			ID:   18,
			Name: "Stop Kontak",
			Url:  "https://ik.imagekit.io/fortyfour/FacilityCategory/electric_bolt_FILL0_wght400_GRAD0_opsz48_6xcn8a2FUk.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670429103611s",
		},
	}

	var count int64
	db.Model(&entity.Category{}).Count(&count)
	if count != 0 {
		return nil
	}

	return db.Create(&category).Error
}

func InitBank(db *gorm.DB) error {
	Bank := entity.Banks{
		{
			ID:   1,
			Name: "BNI",
			Icon: "https://ik.imagekit.io/fortyfour/banks/bni_1__pftwxJOII.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427443148",
		},
		{
			ID:   2,
			Name: "BRI",
			Icon: "https://ik.imagekit.io/fortyfour/banks/bri_dg-04n5YY.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427442469",
		},
		{
			ID:   3,
			Name: "BCA",
			Icon: "https://ik.imagekit.io/fortyfour/banks/bca_YfcLDnPQt.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427476426",
		},
		{
			ID:   4,
			Name: "Mandiri",
			Icon: "https://ik.imagekit.io/fortyfour/banks/mandiri_RHAXhiLalh.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427505115",
		},
		{
			ID:   5,
			Name: "Permata",
			Icon: "https://ik.imagekit.io/fortyfour/banks/permata_RSt5x6SZR.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427674886",
		},
		{
			ID:   6,
			Name: "BTN",
			Icon: "https://ik.imagekit.io/fortyfour/banks/btn_oxksrX1Rc.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427674866",
		},
		{
			ID:   7,
			Name: "Danamon",
			Icon: "https://ik.imagekit.io/fortyfour/banks/danamon_UW76q_DGc.svg?ik-sdk-version=javascript-1.4.3&updatedAt=1670427674777",
		},
	}

	var count int64
	db.Model(&entity.Bank{}).Count(&count)
	if count != 0 {
		return nil
	}

	return db.Create(&Bank).Error
}
