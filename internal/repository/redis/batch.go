package redis

import (
	"metertronik/internal/domain/repository"
	"context"
	"metertronik/internal/domain/entity"
	"time"
	"github.com/redis/go-redis/v9"
	"fmt"
	"encoding/json"
)

type RedisBatchRepo struct {
	client *redis.Client
}

func NewRedisBatchRepo(client *redis.Client) repository.RedisBatchRepo {
	return &RedisBatchRepo{
		client: client,
	}
}

func (r *RedisBatchRepo) GetDailyActivityCache(ctx context.Context, deviceID string, date string) (*entity.DailyElectricity, *[]entity.HourlyElectricity, error){
	key := fmt.Sprintf("daily_activity:%s:%s", deviceID, date)

	data, err := r.client.Get(ctx, key).Result()

	if err != nil {	
		if err == redis.Nil {
			return nil, nil, fmt.Errorf("daily activity cache not found")
		}
		return nil, nil, fmt.Errorf("failed to get daily activity cache: %w", err)
	}

	var daily entity.DailyElectricity
	if err := json.Unmarshal([]byte(data), &daily); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal daily activity cache: %w", err)
	}

	var hourly []entity.HourlyElectricity
	if err := json.Unmarshal([]byte(data), &hourly); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal hourly activity cache: %w", err)
	}

	return &daily, &hourly, nil
}

func (r *RedisBatchRepo) SetDailyActivityCache(ctx context.Context, deviceID string, date string, daily *entity.DailyElectricity, hourly *[]entity.HourlyElectricity, ttl time.Duration) error{
	key := fmt.Sprintf("daily_activity:%s:%s", deviceID, date)

	data, err := json.Marshal(daily)
	if err != nil {
		return fmt.Errorf("failed to marshal daily activity cache: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set daily activity cache: %w", err)
	}

	return nil
}
	
func (r *RedisBatchRepo) GetDailyListCache(ctx context.Context, deviceID string, sortBy string, lastDate string) (*[]entity.DailyElectricity, error){
	key := fmt.Sprintf("daily_list:%s:%s:%s", deviceID, sortBy, lastDate)

	data, err := r.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("daily list cache not found")
		}
	}

	var daily []entity.DailyElectricity
	if err := json.Unmarshal([]byte(data), &daily); err != nil {
		return nil, fmt.Errorf("failed to unmarshal daily list cache: %w", err)
	}

	return &daily, nil
}

func (r *RedisBatchRepo) SetDailyListCache(ctx context.Context, deviceID string, sortBy string, lastDate string, dailyList *[]entity.DailyElectricity, ttl time.Duration) error{
	key := fmt.Sprintf("daily_list:%s:%s:%s", deviceID, sortBy, lastDate)

	data, err := json.Marshal(dailyList)
	if err != nil {
		return fmt.Errorf("failed to marshal daily list cache: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set daily list cache: %w", err)
	}

	return nil
}

func (r *RedisBatchRepo) GetDailyRangeCache(ctx context.Context, deviceID string, start string, end string, lastDate string, limit int) (*[]entity.DailyElectricity, error){
	key := fmt.Sprintf("daily_range:%s:%s:%s:%s:%d", deviceID, start, end, lastDate, limit)

	data, err := r.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("daily range cache not found")
		}
	}

	var daily []entity.DailyElectricity

	if err := json.Unmarshal([]byte(data), &daily); err != nil {
		return nil, fmt.Errorf("failed to unmarshal daily range cache: %w", err)
	}

	return &daily, nil
}

func (r *RedisBatchRepo) SetDailyRangeCache(ctx context.Context, deviceID string, start string, end string, lastDate string, limit int, dailyRange *[]entity.DailyElectricity, ttl time.Duration) error{
	key := fmt.Sprintf("daily_range:%s:%s:%s:%s:%d", deviceID, start, end, lastDate, limit)

	data, err := json.Marshal(dailyRange)

	if err != nil {
		return fmt.Errorf("failed to marshal daily range cache: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set daily range cache: %w", err)
	}

	return nil
}
