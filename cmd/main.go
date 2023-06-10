package main

import (
	"context"
	"fmt"

	"github.com/ryohei1216/go-redis-pub-sub/redis"
)

func main() {
	ctx := context.Background()

	redisClient := redis.New()

	err := redisClient.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
