package main

import (
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	conn *redis.Client
}

func NewRedisClient() *Redis {
	conn := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Redis{
		conn: conn,
	}
}
