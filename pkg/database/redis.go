package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v9"
)

func InitRedis(host string, port string, password string, db string) *redis.Client {
	if db == "" {
		db = "0"
	}

	dbInt, err := strconv.Atoi(db)
	if err != nil {
		log.Fatalf("Error converting redis db to int: %v", err)
	}

	address := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbInt,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error senidng ping to redis: %v", err)
	}

	return client
}
