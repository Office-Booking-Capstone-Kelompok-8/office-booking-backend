package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

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

type RedisClient interface {
	Set(ctx context.Context, key string, value string, exp time.Duration) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
}

type RedisClientImpl struct {
	redis *redis.Client
}

func NewRedisClient(redis *redis.Client) RedisClient {
	return &RedisClientImpl{
		redis: redis,
	}
}

func (t *RedisClientImpl) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	return t.redis.Set(ctx, key, value, exp).Err()
}

func (t *RedisClientImpl) Del(ctx context.Context, key string) error {
	return t.redis.Del(ctx, key).Err()
}

func (t *RedisClientImpl) Get(ctx context.Context, key string) (string, error) {
	return t.redis.Get(ctx, key).Result()
}
