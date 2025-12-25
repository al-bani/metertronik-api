package redis

import (
	"context"
	"log"

	"metertronik/internal/domain/repository"
	repoRedis "metertronik/internal/repository/redis/auth"
	"metertronik/pkg/config"

	"github.com/redis/go-redis/v9"
)

func SetupRedisAuth(cfg *config.Config) (repository.UsersRepoRedis, func()) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis Auth is not available: %v. Token caching will be disabled.", err)
		client.Close()
		return nil, func() {}
	}

	log.Println("Redis Auth connected successfully")
	redisAuthRepo := repoRedis.NewUsersRepoRedis(client)

	cleanup := func() {
		client.Close()
	}

	return redisAuthRepo, cleanup
}
