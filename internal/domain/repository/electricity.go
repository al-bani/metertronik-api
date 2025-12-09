package repository

import (
	"context"
	"metertronik/internal/domain/entity"
	"time"
)

type InfluxRepo interface {
	SaveRealTimeElectricity(ctx context.Context, electricity *entity.RealTimeElectricity) error
	GetRealTimeElectricity(ctx context.Context, deviceID string) (*[]entity.RealTimeElectricity, error)
}

type RedisRepo interface {
	SetLatestElectricity(ctx context.Context, deviceID string, electricity *entity.RealTimeElectricity) error
	GetLatestElectricity(ctx context.Context, deviceID string) (*entity.RealTimeElectricity, error)
	DeleteLatestElectricity(ctx context.Context, deviceID string) error
	SaveElectricityHistory(ctx context.Context, deviceID string, electricity *entity.RealTimeElectricity, ttl time.Duration) error
	HasChanged(ctx context.Context, deviceID string, newData *entity.RealTimeElectricity) (bool, *entity.RealTimeElectricity, error)
}

type PostgresRepo interface {
	SaveHourlyElectricity(ctx context.Context, hourlyElectricity *entity.HourlyElectricity) error
	GetHourlyElectricity(ctx context.Context, deviceID string, dataRange int) (*[]entity.HourlyElectricity, error)
	SaveDailyElectricity(ctx context.Context, dailyElectricity *entity.DailyElectricity) error
	GetTarrifs(ctx context.Context) (*entity.Tarrifs, error)
}
