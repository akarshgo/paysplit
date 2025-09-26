package rediscli

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    "", // set via REDIS_PASSWORD if needed
		DB:          0,
		DialTimeout: 5 * time.Second,
	})

	ctx := context.Background()
	if err := Rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis: ping failed: %v", err)
	}
}
