package main

import (
	"github.com/joho/godotenv"
	"log"
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
	//env := config.LoadConfig()
}
