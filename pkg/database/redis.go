package database

import (
	"context"
	"log"

	"metertronik/internal/domain/repository"
	repoRedis "metertronik/internal/repository/redis"
	"metertronik/pkg/config"

	"github.com/redis/go-redis/v9"
)

func SetupRedisRealtime(cfg *config.Config) (repository.RedisRealtimeRepo, func()) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis is not available: %v. Caching will be disabled.", err)
		client.Close()
		return nil, func() {}
	}

	log.Println("Redis connected successfully")
	RedisRealtimeRepo := repoRedis.NewRedisRealtimeRepo(client)

	cleanup := func() {
		client.Close()
	}

	return RedisRealtimeRepo, cleanup
}

func SetupRedisRealtimeBatch(cfg *config.Config) (repository.RedisBatchRepo, func()) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis is not available: %v. Caching will be disabled.", err)
		client.Close()
		return nil, func() {}
	}

	log.Println("Redis Batch connected successfully")
	RedisBatchRepo := repoRedis.NewRedisBatchRepo(client)

	cleanup := func() {
		client.Close()
	}

	return RedisBatchRepo, cleanup
}
