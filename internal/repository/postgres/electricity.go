package postgres

import (
	"context"
	"fmt"
	"metertronik/internal/domain/entity"

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
