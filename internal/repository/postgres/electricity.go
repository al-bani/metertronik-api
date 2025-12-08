package postgres

import (
	"database/sql"
	"context"
	"metertronik/internal/domain/entity"
)

type ElectricityRepoPostgres struct {
	db *sql.DB
}

func NewElectricityRepoPostgres(db *sql.DB) *ElectricityRepoPostgres {
	return &ElectricityRepoPostgres{
		db: db,
	}
}

func (r *ElectricityRepoPostgres) SaveHourlyElectricity(ctx context.Context, hourlyElectricity *entity.HourlyElectricity) error {
	return nil
}

