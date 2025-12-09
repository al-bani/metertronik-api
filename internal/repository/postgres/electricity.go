package postgres

import (
	"context"
	"fmt"
	"metertronik/internal/domain/entity"
	"metertronik/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type ElectricityRepoPostgres struct {
	db *gorm.DB
}

func NewElectricityRepoPostgres(db *gorm.DB) *ElectricityRepoPostgres {
	return &ElectricityRepoPostgres{
		db: db,
	}
}

func (r *ElectricityRepoPostgres) SaveHourlyElectricity(ctx context.Context, hourlyElectricity *entity.HourlyElectricity) error {
	if err := r.db.WithContext(ctx).Table("hourly_data").Create(hourlyElectricity).Error; err != nil {
		return fmt.Errorf("failed to insert hourly electricity data: %w", err)
	}

	return nil
}

func (r *ElectricityRepoPostgres) GetHourlyElectricity(ctx context.Context, deviceID string, hours int) (*[]entity.HourlyElectricity, error) {
	var hourlyElectricity []entity.HourlyElectricity

	endTime := utils.TimeNowHourly()                                  
	startTime := utils.NewTimeData(endTime.Time.Add(-time.Duration(hours) * time.Hour))

	if err := r.db.WithContext(ctx).Table("hourly_data").Where("device_id = ? AND ts BETWEEN ? AND ?", deviceID, startTime, endTime).Find(&hourlyElectricity).Error; err != nil {
		return nil, fmt.Errorf("failed to get hourly electricity data: %w", err)
	}
	
	return &hourlyElectricity, nil
}

func (r *ElectricityRepoPostgres) SaveDailyElectricity(ctx context.Context, dailyElectricity *entity.DailyElectricity) error {
	if err := r.db.WithContext(ctx).Table("daily_data").Create(dailyElectricity).Error; err != nil {
		return fmt.Errorf("failed to insert daily electricity data: %w", err)
	}
	
	return nil
}


func (r *ElectricityRepoPostgres) GetTarrifs(ctx context.Context) (*entity.Tarrifs, error) {
	var tarrifs entity.Tarrifs

	if err := r.db.WithContext(ctx).Table("tarrifs").Where("effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)", utils.TimeNow(), utils.TimeNow()).First(&tarrifs).Error; err != nil {
		return nil, fmt.Errorf("failed to get tarrifs: %w", err)
	}

	return &tarrifs, nil
}