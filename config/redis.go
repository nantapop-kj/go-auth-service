package config

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func InitRedis() *redis.Client {
	dbStr := GetEnv("REDIS_DB", "0")
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		fmt.Printf("Invalid REDIS_DB value '%s', using 0\n", dbStr)
		db = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_ADDRESS", ""),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       db,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	log.Println("\033[32m✅ [Redis] Connected successfully.\033[0m")

	return rdb
}
