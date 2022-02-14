package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	RDB *redis.Client
)

func init() {
	host := flag.String("host", "localhost", "Redis host")
	port := flag.String("port", "6379", "Redis port")
	flag.Parse()
	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", *host, *port),
	})
	if _, err := RDB.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
}
