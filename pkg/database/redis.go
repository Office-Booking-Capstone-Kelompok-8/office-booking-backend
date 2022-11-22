package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
)

func InitRedis(host string, port string, password string, db string) *redis.Client {
	if db == "" {
		db = "0"
	}

	dbInt, err := strconv.Atoi(db)
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbInt,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return client
}
