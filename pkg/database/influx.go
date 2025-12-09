package database

import (
	"context"
	"log"

	"metertronik/internal/domain/repository"
	repoInflux "metertronik/internal/repository/influx"
	"metertronik/pkg/config"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func SetupInfluxDB(cfg *config.Config) (repository.InfluxRepo, func()) {
	client := influxdb2.NewClient(cfg.InfluxURL, cfg.InfluxToken)

	if _, err := client.Health(context.Background()); err != nil {
		log.Fatalf("InfluxDB health check failed: %v", err)
	}

	repo := repoInflux.NewElectricityRepo(client, cfg.InfluxOrg, cfg.InfluxBucket)
	cleanup := func() {
		client.Close()
	}

	return repo, cleanup
}
