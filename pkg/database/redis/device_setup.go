package redis

import (
	"context"
	"log"

	"metertronik/internal/domain/repository"
	repoRedis "metertronik/internal/repository/redis/api"
	"metertronik/pkg/config"

	"github.com/redis/go-redis/v9"
)

func SetupRedisDevice(cfg *config.Config) (repository.RedisDeviceRepo, func()) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis Device is not available: %v. Device pairing will be disabled.", err)
		client.Close()
		return nil, func() {}
	}

	log.Println("Redis Device connected successfully")
	redisDeviceRepo := repoRedis.NewRedisDeviceRepo(client)

	cleanup := func() {
		client.Close()
	}

	return redisDeviceRepo, cleanup
}