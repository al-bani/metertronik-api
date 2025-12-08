package service

import (
	"context"
	"log"
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
	errInflux := s.influxRepo.SaveRealTimeElectricity(ctx, data)

	if errInflux != nil {
		log.Printf("Error saving real time electricity to influx: %v", errInflux)
	} else {
		//log.Println("Saving data to influxDB ")
	}

	changed, _, err := s.RedisRepo.HasChanged(ctx, data.DeviceID, data)
	if err != nil {
		log.Printf("Error comparing cache: %v", err)
	}

	if !changed {
		log.Printf("No change for device %s (skip caching)", data.DeviceID)
		return nil
	}

	//log.Printf("Data changed for device %s", data.DeviceID)

	if err := s.RedisRepo.SetLatestElectricity(ctx, data.DeviceID, data); err != nil {
		log.Printf("❌ Failed saving latest cache: %v", err)
	} else {
		//log.Printf("✅ Updated latest cache for device %s", data.DeviceID)
	}

	if err := s.RedisRepo.SaveElectricityHistory(ctx, data.DeviceID, data, 5*time.Minute); err != nil {
		log.Printf("❌ Failed saving history cache: %v", err)
	} else {
		//log.Printf("🕒 Saved history for device %s", data.DeviceID)
	}

	return nil
}
