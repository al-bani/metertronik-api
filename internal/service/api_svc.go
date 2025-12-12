package service

import (
	"context"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
)

type ApiService struct {
	postgresRepo repository.PostgresRepo
}

func NewApiService(postgresRepo repository.PostgresRepo) *ApiService {
	return &ApiService{
		postgresRepo: postgresRepo,
	}
}

type DailyActivityResponse struct {
	Daily  *entity.DailyElectricity    `json:"daily"`
	Hourly *[]entity.HourlyElectricity `json:"hourly"`
}

func (s *ApiService) DailyActivity(ctx context.Context, deviceID string, dateStr string) (*DailyActivityResponse, error) {
	date, err := utils.ParseDate(dateStr)

	if err != nil {
		return nil, err
	}

	dailyElectricity, hourlyElectricityList, err := s.postgresRepo.GetDailyElectricity(ctx, deviceID, date)
	if err != nil {
		return nil, err
	}

	response := &DailyActivityResponse{
		Daily:  dailyElectricity,
		Hourly: hourlyElectricityList,
	}

	return response, nil
}

func (s *ApiService) DailyList(ctx context.Context, deviceID string, time string, tariff string) (*[]entity.DailyElectricity, error) {
	sortBy := "day asc"

	if time != "" {
		if time == "asc" {
			sortBy = "day asc"
		} else {
			sortBy = "day desc"
		}
	} else if tariff != "" {
		if tariff == "asc" {
			sortBy = "total_cost asc"
		} else {
			sortBy = "total_cost desc"
		}
	}

	dailyElectricityList, err := s.postgresRepo.GetDailyElectricityList(ctx, deviceID, sortBy)

	if err != nil {
		return nil, err
	}

	return dailyElectricityList, nil
}

func (s *ApiService) DailyRange(ctx context.Context, deviceID string, startStr string, endStr string) (*[]entity.DailyElectricity, error) {
	start, err := utils.ParseDate(startStr)
	if err != nil {
		return nil, err
	}

	end, err := utils.ParseDate(endStr)

	if err != nil {
		return nil, err
	}

	dailyElectricityList, err := s.postgresRepo.GetDailyRange(ctx, deviceID, start, end)

	if err != nil {
		return nil, err
	}

	return dailyElectricityList, nil
}
