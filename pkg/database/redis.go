package database

import (
	"context"
	"log"

	"metertronik/internal/domain/repository"
	repoRedis "metertronik/internal/repository/redis"
	"metertronik/pkg/config"

	"github.com/redis/go-redis/v9"
)

// SetupRedis membuat koneksi ke Redis dan mengembalikan cache repository
// Jika Redis tidak tersedia, akan mengembalikan nil dan cleanup function kosong
func SetupRedis(cfg *config.Config) (repository.RedisRepo, func()) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis tidak tersedia: %v. Caching akan dinonaktifkan.", err)
		client.Close()
		return nil, func() {}
	}

	log.Println("✅ Redis connected successfully")
	RedisRepo := repoRedis.NewRedisRepo(client)

	cleanup := func() {
		client.Close()
	}

	return RedisRepo, cleanup
}
