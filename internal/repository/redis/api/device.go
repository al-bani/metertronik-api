package redis

import (
	"metertronik/internal/domain/repository"
	"github.com/redis/go-redis/v9"
	"context"
	"fmt"
	"metertronik/pkg/utils"
	"log"
)

type RedisDeviceRepo struct {
	client *redis.Client
}

func NewRedisDeviceRepo(client *redis.Client) repository.RedisDeviceRepo {
	return &RedisDeviceRepo{
		client: client,
	}
}

func (r *RedisDeviceRepo) GetDevicePairing(ctx context.Context, pairing string, deviceID string) (int64, error) {
	key := fmt.Sprintf("device:%s:pairing:%s", deviceID, pairing)

	log.Println("Redis, GetDevicePairing : ", key)

	data, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}

	return data, nil
}

func (r *RedisDeviceRepo) SetDevicePairing(ctx context.Context, pairing string, deviceID string, userId int64) error {
	key := fmt.Sprintf("device:%s:pairing:%s", deviceID, pairing)
	ttl := utils.Minutes(10)
	log.Println("Redis, SetDevicePairing : ", key, userId, ttl)

	if err := r.client.Set(ctx, key, userId, ttl).Err(); err != nil {
		return err
	}

	return nil
}