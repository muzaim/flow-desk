// config/redis.go
package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Gagal terhubung ke Redis: %v", err)
	}

	fmt.Println("Koneksi ke Redis berhasil!")
	return rdb
}
