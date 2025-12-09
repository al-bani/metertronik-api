package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
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

func (s *CronService) HourlyAggregation(ctx context.Context) (*entity.HourlyElectricity, error) {
	deviceID := "esp32-meter-001"
	realtimeDataList, err := s.influxRepo.GetRealTimeElectricity(ctx, deviceID)
	if err != nil {
		return &entity.HourlyElectricity{}, err
	}

	if realtimeDataList == nil || len(*realtimeDataList) == 0 {
		log.Printf("No data found for device %s", deviceID)
		return &entity.HourlyElectricity{}, errors.New("no data found for device")
	}

	dataList := *realtimeDataList
	count := len(dataList)

	var totalVoltage, totalCurrent, totalPower, totalFrequency float64
	minPower := dataList[0].Power
	maxPower := dataList[0].Power

	for _, data := range dataList {
		totalVoltage += data.Voltage
		totalCurrent += data.Current
		totalPower += data.Power
		totalFrequency += data.Frequency

		if data.Power < minPower {
			minPower = data.Power
		}
		if data.Power > maxPower {
			maxPower = data.Power
		}
	}

	usageKWh := math.Max(0, dataList[count-1].TotalEnergy-dataList[0].TotalEnergy)

	hourlyData := entity.HourlyElectricity{
		DeviceID:   deviceID,
		UsageKWh:   usageKWh,
		TotalCost:  0,
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
