package redis

import (
	"context"
	"metertronik/internal/domain/repository"
	"fmt"
	"github.com/redis/go-redis/v9"
	"metertronik/pkg/utils"
)

type UsersRepoRedis struct {
	client *redis.Client
}

func NewUsersRepoRedis(client *redis.Client) repository.UsersRepoRedis {
	return &UsersRepoRedis{
		client: client,
	}
}

func (r *UsersRepoRedis) SetToken(ctx context.Context, identifier int64, token string) error {
	key := fmt.Sprintf("auth:refresh:token:%d", identifier)
	ttl := utils.Days(30)

	if err := r.client.Set(ctx, key, token, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (r *UsersRepoRedis) TokenValidation(ctx context.Context, identifier string, token string) (bool, error) {
	// key := fmt.Sprintf("auth:token:%s", identifier)

	// storedToken, err := r.client.Get(ctx, key).Result()
	// if err == redis.Nil {
	// 	return false, nil // Token not found
	// }
	// if err != nil {
	// 	return false, fmt.Errorf("failed to get token: %w", err)
	// }

	// return storedToken == token, nil

	return false, nil
}

func (r *UsersRepoRedis) ResetExpired(ctx context.Context, identifier int64, token string) error {
	key := fmt.Sprintf("auth:refresh:token:%d", identifier)
	ttl := utils.Days(30)

	// ttlLog, _ := r.client.TTL(ctx, key).Result()
	// log.Println("ttl sebelum reset", ttlLog)

	stored, err := r.client.Get(ctx, key).Result()	
	if err != nil || stored != token {
		return fmt.Errorf("failed to reset expired token: %w", err)
	}

	if err := r.client.Set(ctx, key, stored, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	// ttlLog, _ = r.client.TTL(ctx, key).Result()
	// log.Println("ttl setelah reset", ttlLog)

	return nil
}
