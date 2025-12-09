package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	InfluxURL    string
	InfluxToken  string
	InfluxOrg    string
	InfluxBucket string

	RedisAddr     string
	RedisPassword string

	RabbitMQURL string

	Port string

	PGHOST     string
	PGPORT     string
	PGUSER     string
	PGPASSWORD string
	PGDATABASE string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	return &Config{
		InfluxURL:    getEnv("INFLUX_URL", ""),
		InfluxToken:  getEnv("INFLUX_TOKEN", ""),
		InfluxOrg:    getEnv("INFLUX_ORG", ""),
		InfluxBucket: getEnv("INFLUX_BUCKET", ""),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		RabbitMQURL: getEnv("RABBITMQ_URL", ""),

		Port: getEnv("PORT", "8080"),

		PGHOST:     getEnv("PG_HOST", ""),
		PGPORT:     getEnv("PG_PORT", ""),
		PGUSER:     getEnv("PG_USER", ""),
		PGPASSWORD: getEnv("PG_PASS", ""),
		PGDATABASE: getEnv("PG_DTBS", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
