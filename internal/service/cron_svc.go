package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (s *CronService) DailyAggregation(ctx context.Context) (*entity.DailyElectricity, error) {
	deviceID := "device-001"
	hourlyDataList, err := s.postgresRepo.GetHourlyElectricity(ctx, deviceID, 24, nil)

	if err != nil {
		return nil, err
	}

	if hourlyDataList == nil || len(*hourlyDataList) == 0 {
		log.Printf("No data found for device %s", deviceID)
		return nil, errors.New("no data found for device")
	}

	dataList := *hourlyDataList
	count := len(dataList)

	var totalVoltage, totalCurrent, totalPower, Energy float64

	minPower := dataList[0].MinPower
	maxPower := dataList[0].MaxPower

	for _, data := range dataList {
		totalVoltage += data.AvgVoltage
		totalCurrent += data.AvgCurrent
		totalPower += data.AvgPower
		Energy += data.Energy

		if data.MinPower < minPower {
			minPower = data.MinPower
		}
		if data.MaxPower > maxPower {
			maxPower = data.MaxPower
		}
	}

	tarrifs, err := s.postgresRepo.GetTarrifs(ctx)

	if err != nil {
		log.Printf("Failed to get tarrifs: %v", err)
		return nil, fmt.Errorf("failed to get tarrifs: %w", err)
	}

	dailyData := entity.DailyElectricity{
		DeviceID:   deviceID,
		Energy:     Energy,
		TotalCost:  (Energy * tarrifs.PricePerKwh) * 1.10,
		AvgVoltage: totalVoltage / float64(count),
		AvgCurrent: totalCurrent / float64(count),
		AvgPower:   totalPower / float64(count),
		MinPower:   minPower,
		MaxPower:   maxPower,
		Day:        utils.TimeNowDaily(),
		CreatedAt:  utils.TimeNow(),
	}

	log.Println("Daily Data \ndeviceID: ", dailyData.DeviceID, "\nEnergy: ", dailyData.Energy, "\ntotalCost: ", dailyData.TotalCost, "\navgVoltage: ", dailyData.AvgVoltage, "\navgCurrent: ", dailyData.AvgCurrent, "\navgPower: ", dailyData.AvgPower, "\nminPower: ", dailyData.MinPower, "\nmaxPower: ", dailyData.MaxPower, "\nDay: ", dailyData.Day, "\nCreatedAt: ", dailyData.CreatedAt)

	if err := s.postgresRepo.SaveDailyElectricity(ctx, &dailyData); err != nil {
		log.Printf("Failed to save daily data to postgres: %v", err)
		return nil, fmt.Errorf("failed to save daily data to postgres: %w", err)
	}

	return nil, nil
}

func (s *CronService) HourlyAggregation(ctx context.Context) (*entity.HourlyElectricity, error) {
	deviceID := "device-001"
	realtimeDataList, err := s.influxRepo.GetRealTimeElectricity(ctx, deviceID)
	if err != nil {
		return &entity.HourlyElectricity{}, err
	}

	if realtimeDataList == nil || len(*realtimeDataList) == 0 {
		log.Printf("No data found for device %s", deviceID)
		return &entity.HourlyElectricity{}, errors.New("no data found for device")
	}

	tarrifs, err := s.postgresRepo.GetTarrifs(ctx)

	if err != nil {
		log.Printf("Failed to get tarrifs: %v", err)
		return &entity.HourlyElectricity{}, fmt.Errorf("failed to get tarrifs: %w", err)
	}

	dataList := *realtimeDataList
	count := len(dataList)

	var totalVoltage, totalCurrent, totalPower, totalFrequency, Energy float64
	minPower := dataList[0].Power
	maxPower := dataList[0].Power

	for _, data := range dataList {
		totalVoltage += data.Voltage
		totalCurrent += data.Current
		totalPower += data.Power
		totalFrequency += data.Frequency
		Energy += data.Energy

		if data.Power < minPower {
			minPower = data.Power
		}
		if data.Power > maxPower {
			maxPower = data.Power
		}
	}

	hourlyData := entity.HourlyElectricity{
		DeviceID:   deviceID,
		Energy:     Energy,
		TotalCost:  (Energy * tarrifs.PricePerKwh) * 1.10,
		AvgVoltage: totalVoltage / float64(count),
		AvgCurrent: totalCurrent / float64(count),
		AvgPower:   totalPower / float64(count),
		MinPower:   minPower,
		MaxPower:   maxPower,
		TS:         utils.TimeNowHourly(),
		CreatedAt:  utils.TimeNow(),
	}

	if err := s.postgresRepo.SaveHourlyElectricity(ctx, &hourlyData); err != nil {
		log.Printf("Failed to save hourly data to postgres: %v", err)
		return &hourlyData, fmt.Errorf("failed to save hourly data to postgres: %w", err)
	}

	return &hourlyData, nil
}
