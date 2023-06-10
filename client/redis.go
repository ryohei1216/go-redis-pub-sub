package redis

import go_redis "github.com/go-redis/redis/v8"

func New() *go_redis.Client {
	client := go_redis.NewClient(&go_redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}
