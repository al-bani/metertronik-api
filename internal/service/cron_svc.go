package service

import (
	"context"
	"log"
	"math"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
	"errors"
)

type CronService struct {
	influxRepo repository.InfluxRepo
}

func NewCronService(influxRepo repository.InfluxRepo) *CronService {
	return &CronService{
		influxRepo: influxRepo,
	}
}

func (s *CronService) HourlyAggregation(ctx context.Context) (*entity.HourlyElectricity, error) {
	deviceID := "esp32-meter-001"
	realtimeDataList, err := s.influxRepo.GetRealTimeElectricity(ctx, deviceID)
	if err != nil {
		return &entity.HourlyElectricity{}, err
	}

	if realtimeDataList == nil || len(*realtimeDataList) == 0 {
		log.Printf("Tidak ada data untuk device %s", deviceID)
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
		DeviceID:     deviceID,
		UsageKWh:     usageKWh,
		TotalCost:    0,
		AvgVoltage:   totalVoltage / float64(count),
		AvgCurrent:   totalCurrent / float64(count),
		AvgPower:     totalPower / float64(count),
		AvgFrequency: totalFrequency / float64(count),
		MinPower:     minPower,
		MaxPower:     maxPower,
		CreatedAt:    utils.TimeNow(),
	}

	return &hourlyData, nil
}
