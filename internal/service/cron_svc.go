package service

import (
	"context"
	"errors"
	"time"

	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
)

type CronService struct {
	influxRepo   repository.InfluxRepo
	postgresRepo repository.PostgresRepo
}

func NewCronService(influxRepo repository.InfluxRepo, postgresRepo repository.PostgresRepo) *CronService {
	return &CronService{
		influxRepo:   influxRepo,
		postgresRepo: postgresRepo,
	}
}

func (s *CronService) HourlyAggregation(
	ctx context.Context,
	targetHour time.Time,
	deviceID string,
) (*entity.HourlyElectricity, error) {

	start := utils.NewTimeData(targetHour)
	end := utils.NewTimeData(targetHour.Add(time.Hour))

	realtimeDataList, err := s.influxRepo.
		GetRealTimeElectricityRange(ctx, deviceID, start, end)
	if err != nil || realtimeDataList == nil {
		return nil, errors.New("no realtime data for hour")
	}

	tarrifs, err := s.postgresRepo.GetTarrifs(ctx)
	if err != nil {
		return nil, err
	}

	dataList := *realtimeDataList
	count := len(dataList)

	var totalVoltage, totalCurrent, totalPower, totalFrequency, energy float64
	minPower := dataList[0].Power
	maxPower := dataList[0].Power

	for _, d := range dataList {
		totalVoltage += d.Voltage
		totalCurrent += d.Current
		totalPower += d.Power
		totalFrequency += d.Frequency
		energy += d.Energy

		if d.Power < minPower {
			minPower = d.Power
		}
		if d.Power > maxPower {
			maxPower = d.Power
		}
	}

	hourly := entity.HourlyElectricity{
		DeviceID:   deviceID,
		Energy:     energy,
		TotalCost:  (energy * tarrifs.PricePerKwh) * 1.10,
		AvgVoltage: totalVoltage / float64(count),
		AvgCurrent: totalCurrent / float64(count),
		AvgPower:   totalPower / float64(count),
		MinPower:   minPower,
		MaxPower:   maxPower,
		TS:         start,
		CreatedAt:  utils.TimeNow(),
	}

	return &hourly, s.postgresRepo.UpsertHourlyElectricity(ctx, &hourly)
}

func (s *CronService) DailyAggregation(
	ctx context.Context,
	targetDay time.Time,
	deviceID string,
) (*entity.DailyElectricity, error) {

	start := utils.NewTimeData(targetDay)
	end := utils.NewTimeData(targetDay.AddDate(0, 0, 1))

	hourlyDataList, err := s.postgresRepo.
		GetHourlyElectricityRange(ctx, deviceID, start, end)
	if err != nil || hourlyDataList == nil {
		return nil, errors.New("no hourly data for day")
	}

	dataList := *hourlyDataList
	count := len(dataList)

	var totalVoltage, totalCurrent, totalPower, energy float64
	minPower := dataList[0].MinPower
	maxPower := dataList[0].MaxPower

	for _, d := range dataList {
		totalVoltage += d.AvgVoltage
		totalCurrent += d.AvgCurrent
		totalPower += d.AvgPower
		energy += d.Energy

		if d.MinPower < minPower {
			minPower = d.MinPower
		}
		if d.MaxPower > maxPower {
			maxPower = d.MaxPower
		}
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
		Day:        start,
		CreatedAt:  utils.TimeNow(),
	}

	return &daily, s.postgresRepo.UpsertDailyElectricity(ctx, &daily)
}

func (s *CronService) MonthlyAggregation(
	ctx context.Context,
	targetMonth utils.TimeData,
	deviceID string,
) (*entity.MonthlyElectricity, error) {

	var dailyList *[]entity.DailyElectricity
	var err error

	if targetMonth.IsFirstDayOfMonth() {
		dailyElectricity, _, err := s.postgresRepo.GetDailyElectricity(ctx, deviceID, targetMonth)
		if err != nil {
			return nil, err
		}

		if dailyElectricity != nil {
			dailyList = &[]entity.DailyElectricity{*dailyElectricity}
		} else {
			return nil, errors.New("no daily data for month")
		}
	} else {
		startDate, endDate := targetMonth.GetMonthlyRangeDates()

		dailyList, err = s.postgresRepo.GetDailyRange(ctx, deviceID, startDate, endDate, nil, 32)
		if err != nil {
			return nil, err
		}
	}

	if dailyList == nil || len(*dailyList) == 0 {
		return nil, errors.New("no daily data for month")
	}

	dataList := *dailyList
	var totalEnergy, totalCost float64

	for _, d := range dataList {
		totalEnergy += d.Energy
		totalCost += d.TotalCost
	}

	monthly := entity.MonthlyElectricity{
		DeviceID:  deviceID,
		Month:     targetMonth.StartOfMonth(),
		Energy:    totalEnergy,
		TotalCost: totalCost,
		CreatedAt: utils.TimeNow(),
	}

	return &monthly, s.postgresRepo.UpsertMonthlyElectricity(ctx, &monthly)
}
