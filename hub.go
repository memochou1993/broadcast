package main

import (
	"github.com/go-redis/redis/v8"
)

type Hub struct {
	rdb *redis.Client
}

func newHub() *Hub {
	return &Hub{
		rdb: NewRDB(),
	}
}
