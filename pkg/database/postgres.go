package database

import (
	"fmt"
	"log"

	"metertronik/internal/domain/repository"
	repoPostgres "metertronik/internal/repository/postgres"
	"metertronik/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupPostgres(cfg *config.Config) (repository.PostgresRepo, func()) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGHOST, cfg.PGPORT, cfg.PGUSER, cfg.PGPASSWORD, cfg.PGDATABASE, cfg.PGSSLMODE,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	repo := repoPostgres.NewElectricityRepoPostgres(db)

	cleanup := func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting sql.DB from GORM: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}

	return repo, cleanup
}
