package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) repository.RedisRepo {
	return &RedisRepo{
		client: client,
	}
}

func isSame(a, b *entity.RealTimeElectricity) bool {
	if a == nil || b == nil {
		return false
	}

	return a.Voltage == b.Voltage &&
		a.Current == b.Current &&
		a.Power == b.Power &&
		a.Energy == b.Energy &&
		a.PowerFactor == b.PowerFactor &&
		a.Frequency == b.Frequency
}

func (r *RedisRepo) SetLatestElectricity(ctx context.Context, deviceID string, electricity *entity.RealTimeElectricity) error {
	key := fmt.Sprintf("electricity:latest:%s", deviceID)

	data, err := json.Marshal(electricity)
	if err != nil {
		return fmt.Errorf("failed to marshal latest electricity: %w", err)
	}

	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to set latest cache: %w", err)
	}

	return nil
}

func (r *RedisRepo) GetLatestElectricity(ctx context.Context, deviceID string) (*entity.RealTimeElectricity, error) {
	key := fmt.Sprintf("electricity:latest:%s", deviceID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest cache: %w", err)
	}

	var electricity entity.RealTimeElectricity
	if err := json.Unmarshal([]byte(data), &electricity); err != nil {
		return nil, fmt.Errorf("failed to unmarshal latest electricity: %w", err)
	}

	return &electricity, nil
}

func (r *RedisRepo) SaveElectricityHistory(ctx context.Context, deviceID string, electricity *entity.RealTimeElectricity, ttl time.Duration) error {
	key := fmt.Sprintf("electricity:%s:%s",
		deviceID,
		electricity.CreatedAt.Format(),
	)

	data, err := json.Marshal(electricity)
	if err != nil {
		return fmt.Errorf("failed to marshal history electricity: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set history cache: %w", err)
	}

	return nil
}

func (r *RedisRepo) HasChanged(ctx context.Context, deviceID string, newData *entity.RealTimeElectricity) (bool, *entity.RealTimeElectricity, error) {
	oldData, err := r.GetLatestElectricity(ctx, deviceID)

	if err != nil {
		return false, nil, err
	}

	if oldData == nil {
		return true, nil, nil
	}

	same := isSame(oldData, newData)
	return !same, oldData, nil
}

func (r *RedisRepo) DeleteLatestElectricity(ctx context.Context, deviceID string) error {
	key := fmt.Sprintf("electricity:latest:%s", deviceID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}
