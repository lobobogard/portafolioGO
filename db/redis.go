package db

import "github.com/go-redis/redis/v8"

func Rdb() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return rdb
}
