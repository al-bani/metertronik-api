package service

import (
	"context"
	"errors"
	"log"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
)

type ApiService struct {
	postgresRepo   repository.PostgresRepo
	redisBatchRepo repository.RedisBatchRepo
}

func NewApiService(postgresRepo repository.PostgresRepo, redisBatchRepo repository.RedisBatchRepo) *ApiService {
	return &ApiService{
		postgresRepo:   postgresRepo,
		redisBatchRepo: redisBatchRepo,
	}
}


type DailyActivityResponse struct {
	Daily  *entity.DailyElectricity    `json:"daily"`
	Hourly *[]entity.HourlyElectricity `json:"hourly"`
}

type MonthlyResponse struct {
	Month   *entity.MonthlyElectricity   `json:"month"`
	Daily   *[]entity.DailyElectricity   `json:"daily"`
	Monthly *[]entity.MonthlyElectricity `json:"monthly"`
}

func (s *ApiService) DailyActivity(ctx context.Context, deviceID string, dateStr string) (*DailyActivityResponse, error) {
	date, err := utils.ParseDate(dateStr)

	var dailyElectricity *entity.DailyElectricity
	var hourlyElectricityList *[]entity.HourlyElectricity

	if err != nil {
		return nil, err
	}

	if date.Time.IsZero() {
		return nil, errors.New("date parameter is required")
	}

	if s.redisBatchRepo != nil {
		dailyElectricity, hourlyElectricityList, err = s.redisBatchRepo.GetDailyActivityCache(ctx, deviceID, dateStr)
		if err == nil {
			response := &DailyActivityResponse{
				Daily:  dailyElectricity,
				Hourly: hourlyElectricityList,
			}
			return response, nil
		}
	}

	dailyElectricity, hourlyElectricityList, err = s.postgresRepo.GetDailyElectricity(ctx, deviceID, date)
	if err != nil {
		return nil, err
	}

	if s.redisBatchRepo != nil {
		duration := utils.Days(0)

		if date == utils.TimeNowDaily() {
			duration = utils.Minutes(5)
		} else {
			duration = utils.Days(30)
		}

		err = s.redisBatchRepo.SetDailyActivityCache(ctx, deviceID, dateStr, dailyElectricity, hourlyElectricityList, duration)
		if err != nil {
			return nil, err
		}
	}

	response := &DailyActivityResponse{
		Daily:  dailyElectricity,
		Hourly: hourlyElectricityList,
	}

	return response, nil
}

func (s *ApiService) DailyList(ctx context.Context, deviceID string, time string, tariff string, last string) (*[]entity.DailyElectricity, error) {
	sortBy := "day desc"

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

	var lastDate *utils.TimeData

	if last != "" {
		lastDateData, err := utils.ParseDate(last)
		if err != nil {
			return nil, err
		}
		lastDate = &lastDateData
	}

	var dailyElectricityList *[]entity.DailyElectricity
	var err error

	if s.redisBatchRepo != nil {
		dailyElectricityList, err = s.redisBatchRepo.GetDailyListCache(ctx, deviceID, sortBy, last)
		log.Println("Getting List from Cache")
		if err == nil {
			return dailyElectricityList, nil
		}
	}

	log.Println("Getting List from Postgres")
	dailyElectricityList, err = s.postgresRepo.GetDailyElectricityList(ctx, deviceID, sortBy, lastDate)
	if err != nil {
		return nil, err
	}

	if s.redisBatchRepo != nil {
		log.Println("Setting List to Cache")
		err = s.redisBatchRepo.SetDailyListCache(ctx, deviceID, sortBy, last, dailyElectricityList, utils.Minutes(2))
		if err != nil {
			return nil, err
		}
	}

	return dailyElectricityList, nil
}

func (s *ApiService) DailyRange(ctx context.Context, deviceID string, startStr string, endStr string, last string, limit int) (*[]entity.DailyElectricity, error) {
	start, err := utils.ParseDate(startStr)
	if err != nil {
		return nil, err
	}

	if start.Time.IsZero() {
		return nil, errors.New("start date parameter is required")
	}

	end, err := utils.ParseDate(endStr)
	if err != nil {
		return nil, err
	}

	if end.Time.IsZero() {
		return nil, errors.New("end date parameter is required")
	}

	var lastDate *utils.TimeData

	if last != "" {
		lastDateData, err := utils.ParseDate(last)
		if err != nil {
			return nil, err
		}
		lastDate = &lastDateData
	}

	var dailyElectricityList *[]entity.DailyElectricity

	if s.redisBatchRepo != nil {
		dailyElectricityList, err = s.redisBatchRepo.GetDailyRangeCache(ctx, deviceID, startStr, endStr, last, limit)
		if err == nil {
			return dailyElectricityList, nil
		}
	}

	dailyElectricityList, err = s.postgresRepo.GetDailyRange(ctx, deviceID, start, end, lastDate, limit)
	if err != nil {
		return nil, err
	}

	if s.redisBatchRepo != nil {
		err = s.redisBatchRepo.SetDailyRangeCache(ctx, deviceID, startStr, endStr, last, limit, dailyElectricityList, utils.Days(30))
		if err != nil {
			return nil, err
		}
	}

	return dailyElectricityList, nil
}

func (s *ApiService) MonthlyList(ctx context.Context, deviceID string, dateStr string) (*MonthlyResponse, error) {
	date, err := utils.ParseDate(dateStr)
	if err != nil {
		return nil, err
	}

	var dailyList *[]entity.DailyElectricity

	if date.IsFirstDayOfMonth() {

		dailyActivity, err := s.DailyActivity(ctx, deviceID, dateStr)
		if err != nil {
			return nil, err
		}

		if dailyActivity.Daily != nil {
			dailyList = &[]entity.DailyElectricity{*dailyActivity.Daily}
		} else {
			dailyList = &[]entity.DailyElectricity{}
		}
	} else {
		startDate, endDate := date.GetMonthlyRangeDates()

		startStr := startDate.Format()
		endStr := endDate.Format()

		dailyList, err = s.DailyRange(ctx, deviceID, startStr, endStr, "", 32)
		if err != nil {
			return nil, err
		}
	}

	monthlyList, err := s.postgresRepo.GetMonthlyElectricity(ctx, deviceID)

	if err != nil {
		return nil, err
	}

	currentMonth := date.StartOfMonth()
	currentMonthUTC := currentMonth.Time.UTC()

	var monthly *entity.MonthlyElectricity
	var monthlyListFiltered *[]entity.MonthlyElectricity


	if monthlyList != nil && len(*monthlyList) > 0 {
		var currentMonthData *entity.MonthlyElectricity
		var otherMonths []entity.MonthlyElectricity

		for i := range *monthlyList {
			monthData := &(*monthlyList)[i]
			monthStart := monthData.Month.StartOfMonth()
			monthStartUTC := monthStart.Time.UTC()

			isCurrentMonth := monthStartUTC.Year() == currentMonthUTC.Year() &&
				monthStartUTC.Month() == currentMonthUTC.Month()
				
			if isCurrentMonth {
				currentMonthData = monthData
			} else {
				otherMonths = append(otherMonths, *monthData)
			}
		}

		monthly = currentMonthData

		if len(otherMonths) > 0 {
			monthlyListFiltered = &otherMonths
		} else {
			monthlyListFiltered = &[]entity.MonthlyElectricity{}
		}
	} else {
		monthlyListFiltered = &[]entity.MonthlyElectricity{}
	}

	return &MonthlyResponse{
		Month:   monthly,
		Daily:   dailyList,
		Monthly: monthlyListFiltered,
	}, nil
}

func (s *ApiService) DayNowActivity(ctx context.Context, deviceID string) (*DailyActivityResponse, error) {
	endTime := utils.TimeNowHourly()
	startTime := endTime.StartOfDay()

	hourlyDataList, err := s.postgresRepo.GetHourlyElectricityRange(ctx, deviceID, startTime, endTime)

	count := len(*hourlyDataList)

	var totalVoltage, totalCurrent, totalPower, energy float64
	minPower := (*hourlyDataList)[0].MinPower
	maxPower := (*hourlyDataList)[0].MaxPower

	for _, d := range *hourlyDataList {
		totalVoltage += d.AvgVoltage
		totalCurrent += d.AvgCurrent
		totalPower += d.AvgPower
		energy += d.Energy
	}

	tarrifs, err := s.postgresRepo.GetTarrifs(ctx)
	if err != nil {
		return nil, err
	}

	daily := entity.DailyElectricity{
		DeviceID:   deviceID,
		Energy:     energy,
		TotalCost:  (energy * tarrifs.PricePerKwh) * 1.10,
		AvgVoltage: totalVoltage / float64(count),
		AvgCurrent: totalCurrent / float64(count),
		AvgPower:   totalPower / float64(count),
		MinPower:   minPower,
		MaxPower:   maxPower,
		Day:        startTime,
		CreatedAt:  utils.TimeNow(),
	}

	return &DailyActivityResponse{
		Daily:  &daily,
		Hourly: hourlyDataList,
	}, nil
}