package redis

import (
	"context"
	"metertronik/internal/domain/repository"
	"fmt"
	"github.com/redis/go-redis/v9"
	"metertronik/pkg/utils"
	"log"
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

func (r *UsersRepoRedis) ResetExpired(ctx context.Context, identifier int64, token string) error {
	key := fmt.Sprintf("auth:refresh:token:%d", identifier)
	ttl := utils.Days(30)

	ttlLog, _ := r.client.TTL(ctx, key).Result()
	log.Println("ttl sebelum reset", ttlLog)

	stored, err := r.client.Get(ctx, key).Result()	
	if err != nil || stored != token {
		return fmt.Errorf("failed to reset expired token: %w", err)
	}

	if err := r.client.Set(ctx, key, stored, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	ttlLog, _ = r.client.TTL(ctx, key).Result()
	log.Println("ttl setelah reset", ttlLog)

	return nil
}

func (r *UsersRepoRedis) SetVerificationCodeOtp(ctx context.Context, identifier string, hashCode string) error {
	key := fmt.Sprintf("auth:verification:email:code:%s", identifier)
	ttl := utils.Minutes(5)

	if err := r.client.Set(ctx, key, hashCode, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save verification email code: %w", err)
	}

	return nil
}

func (r *UsersRepoRedis) GetVerificationCodeOtp(ctx context.Context, identifier string) (string, error) {
	key := fmt.Sprintf("auth:verification:email:code:%s", identifier)

	stored, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get verification email code: %w", err)
	}

	return stored, nil
}

func (r *UsersRepoRedis) DeleteToken(ctx context.Context, identifier int64) error {
	key := fmt.Sprintf("auth:refresh:token:%d", identifier)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}