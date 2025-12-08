package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// InfluxDB
	InfluxURL    string
	InfluxToken  string
	InfluxOrg    string
	InfluxBucket string

	// Redis
	RedisAddr     string
	RedisPassword string

	// RabbitMQ
	RabbitMQURL string

	// Server
	Port string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env") // Ignore error jika file tidak ada

	return &Config{
		InfluxURL:    getEnv("INFLUX_URL", ""),
		InfluxToken:  getEnv("INFLUX_TOKEN", ""),
		InfluxOrg:    getEnv("INFLUX_ORG", ""),
		InfluxBucket: getEnv("INFLUX_BUCKET", ""),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		RabbitMQURL: getEnv("RABBITMQ_URL", ""),

		Port: getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
