package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(addr string, password string, db int) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	// Probar la conexión
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error conectando a Redis: %v", err)
	}

	return client, nil
}
