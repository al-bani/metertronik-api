package service

import (
	"context"
	"errors"
	"log"
	"math"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"time"
)

type IngestService struct {
	influxRepo repository.InfluxRepo
	RedisRepo  repository.RedisRepo
}

func NewIngestService(influxRepo repository.InfluxRepo, RedisRepo repository.RedisRepo) *IngestService {
	return &IngestService{
		influxRepo: influxRepo,
		RedisRepo:  RedisRepo,
	}
}

func (s *IngestService) ProcessRealTimeElectricity(ctx context.Context, data *entity.RealTimeElectricity) error {
	previousData, err := s.RedisRepo.GetLatestElectricity(ctx, data.DeviceID)

	if err != nil {
		log.Printf("Error getting previous electricity data: %v", err)
		return err
	} else if previousData == nil {
		data.PowerSurge = 0
		data.PSPercent = 0
		return errors.New("no previous electricity data found")
	} else {
		data.PowerSurge = math.Abs(data.Power - previousData.Power)

		if previousData.Power != 0 {
			data.PSPercent = math.Abs(((data.Power - previousData.Power) / previousData.Power) * 100)
		} else {
			return errors.New("error calculating PSPercent")
		}
	}

	errInflux := s.influxRepo.SaveRealTimeElectricity(ctx, data)

	if errInflux != nil {
		log.Printf("Error saving real time electricity to influx: %v", errInflux)
	} else {
		log.Println("Saving data to influxDB : ", data)
	}

	changed, _, err := s.RedisRepo.HasChanged(ctx, data.DeviceID, data)
	if err != nil {
		log.Printf("Error comparing cache: %v", err)
	}

	if !changed {
		log.Printf("No change for device %s (skip caching)", data.DeviceID)
		return nil
	}

	if err := s.RedisRepo.SetLatestElectricity(ctx, data.DeviceID, data); err != nil {
		log.Printf("❌ Failed saving latest cache: %v", err)
	} else {
		log.Println("✅ Updated latest cache data :", data)
	}

	if err := s.RedisRepo.SaveElectricityHistory(ctx, data.DeviceID, data, 5*time.Minute); err != nil {
		log.Printf("❌ Failed saving history cache: %v", err)
	}

	return nil
}
