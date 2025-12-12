package postgres

import (
	"context"
	"fmt"
	"metertronik/internal/domain/entity"
	"metertronik/pkg/utils"

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

func (r *ElectricityRepoPostgres) GetHourlyElectricity(ctx context.Context, deviceID string, hours int, date *utils.TimeData) (*[]entity.HourlyElectricity, error) {
	var hourlyElectricity []entity.HourlyElectricity

	var endTime utils.TimeData
	var startTime utils.TimeData

	if date == nil {
		endTime = utils.TimeNowHourly()
		startTime = endTime.AddHours(-hours)
	} else {
		startOfDay := date.StartOfDay()
		endOfDay := date.EndOfDay()
		now := utils.TimeNow()

		if endOfDay.Time.After(now.Time) {
			endTime = utils.TimeNowHourly()
		} else {
			endTime = endOfDay.TruncateHour()
		}

		startTime = endTime.AddHours(-hours)

		if startTime.Time.Before(startOfDay.Time) {
			startTime = startOfDay.TruncateHour()
		}
	}

	if err := r.db.WithContext(ctx).Table("hourly_data").Where("device_id = ? AND ts BETWEEN ? AND ?", deviceID, startTime, endTime).Find(&hourlyElectricity).Error; err != nil {
		return nil, fmt.Errorf("failed to get hourly electricity data: %w", err)
	}

	return &hourlyElectricity, nil
}

func (r *ElectricityRepoPostgres) GetDailyElectricity(ctx context.Context, deviceID string, date utils.TimeData) (*entity.DailyElectricity, *[]entity.HourlyElectricity, error) {
	var dailyElectricity entity.DailyElectricity

	if err := r.db.WithContext(ctx).Table("daily_data").Where("device_id = ? AND day = ?", deviceID, date).First(&dailyElectricity).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get daily electricity data: %w", err)
	}

	hourlyElectricityList, err := r.GetHourlyElectricity(ctx, deviceID, 24, &date)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get hourly electricity data: %w", err)
	}

	return &dailyElectricity, hourlyElectricityList, nil
}

func (r *ElectricityRepoPostgres) GetDailyElectricityList(ctx context.Context, deviceID string, sortBy string) (*[]entity.DailyElectricity, error) {
	var dailyElectricityList []entity.DailyElectricity

	if err := r.db.WithContext(ctx).Table("daily_data").Where("device_id = ?", deviceID).Order(sortBy).Find(&dailyElectricityList).Error; err != nil {
		return nil, fmt.Errorf("failed to get daily electricity data list: %w", err)
	}

	return &dailyElectricityList, nil
}

func (r *ElectricityRepoPostgres) GetDailyRange(ctx context.Context, deviceID string, start utils.TimeData, end utils.TimeData) (*[]entity.DailyElectricity, error) {
	var dailyElectricityList []entity.DailyElectricity

	if err := r.db.WithContext(ctx).Table("daily_data").Where("device_id = ? AND day BETWEEN ? AND ?", deviceID, start, end).Order("day asc").Find(&dailyElectricityList).Error; err != nil {
		return nil, fmt.Errorf("failed to get daily electricity data range: %w", err)
	}

	return &dailyElectricityList, nil
}