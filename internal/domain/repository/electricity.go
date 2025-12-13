package repository

import (
	"context"
	"metertronik/internal/domain/entity"
	"metertronik/pkg/utils"
	"time"
)

type InfluxRepo interface {
	SaveRealTimeElectricity(ctx context.Context, electricity *entity.RealTimeElectricity) error
	GetRealTimeElectricity(ctx context.Context, deviceID string) (*[]entity.RealTimeElectricity, error)
}

type RedisRealtimeRepo interface {
	SetLatestElectricity(ctx context.Context, deviceID string, electricity *entity.RealTimeElectricity) error
	GetLatestElectricity(ctx context.Context, deviceID string) (*entity.RealTimeElectricity, error)
	DeleteLatestElectricity(ctx context.Context, deviceID string) error
	SaveElectricityHistory(ctx context.Context, deviceID string, electricity *entity.RealTimeElectricity, ttl time.Duration) error
	HasChanged(ctx context.Context, deviceID string, newData *entity.RealTimeElectricity) (bool, *entity.RealTimeElectricity, error)
}

type RedisBatchRepo interface {
	GetDailyActivityCache(ctx context.Context, deviceID string, date string) (*entity.DailyElectricity, *[]entity.HourlyElectricity, error)
	SetDailyActivityCache(ctx context.Context, deviceID string, date string, daily *entity.DailyElectricity, hourly *[]entity.HourlyElectricity, ttl time.Duration) error
	
	GetDailyListCache(ctx context.Context, deviceID string, sortBy string, lastDate string) (*[]entity.DailyElectricity, error)
	SetDailyListCache(ctx context.Context, deviceID string, sortBy string, lastDate string, data *[]entity.DailyElectricity, ttl time.Duration) error

	GetDailyRangeCache(ctx context.Context, deviceID string, start string, end string, lastDate string) (*[]entity.DailyElectricity, error)
	SetDailyRangeCache(ctx context.Context, deviceID string, start string, end string, lastDate string, data *[]entity.DailyElectricity, ttl time.Duration) error
}

type PostgresRepo interface {
	SaveHourlyElectricity(ctx context.Context, hourlyElectricity *entity.HourlyElectricity) error
	GetHourlyElectricity(ctx context.Context, deviceID string, hours int, date *utils.TimeData) (*[]entity.HourlyElectricity, error)
	SaveDailyElectricity(ctx context.Context, dailyElectricity *entity.DailyElectricity) error
	GetDailyElectricity(ctx context.Context, deviceID string, date utils.TimeData) (*entity.DailyElectricity, *[]entity.HourlyElectricity, error)
	GetTarrifs(ctx context.Context) (*entity.Tarrifs, error)
	GetDailyElectricityList(ctx context.Context, deviceID string, sortBy string, lastDate *utils.TimeData) (*[]entity.DailyElectricity, error)
	GetDailyRange(ctx context.Context, deviceID string, start utils.TimeData, end utils.TimeData, lastDate *utils.TimeData) (*[]entity.DailyElectricity, error)
}
