package redis

import (
	"context"
	"log"

	go_redis "github.com/go-redis/redis/v8"
)

func New(ctx context.Context) *go_redis.Client {
	client := go_redis.NewClient(&go_redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(ctx).Err()

	if err != nil {
		log.Fatal("failed to connect redis", err)
	}

	return client
}
