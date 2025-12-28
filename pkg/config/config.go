package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	InfluxURL    string
	InfluxToken  string
	InfluxOrg    string
	InfluxBucket string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	RabbitMQURL           string
	RabbitMQQueueName     string
	RabbitMQRoutingKey    string
	RabbitMQExchange      string
	RabbitMQPrefetchCount int
	RabbitMQRetryDelay    time.Duration

	Port    string
	GinMode string

	PGHOST     string
	PGPORT     string
	PGUSER     string
	PGPASSWORD string
	PGDATABASE string
	PGSSLMODE  string

	CORSAllowOrigins []string
	CORSAllowMethods []string
	CORSAllowHeaders []string

	CronHourlyInterval time.Duration
	CronDailyInterval  time.Duration

	ConsumerLogInterval time.Duration

	SendgridAPIKey string
	SendgridFromEmail string
	SendgridFromName string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	rabbitMQPrefetchCount, _ := strconv.Atoi(getEnv("RABBITMQ_PREFETCH_COUNT", "50"))
	rabbitMQRetryDelaySeconds, _ := strconv.Atoi(getEnv("RABBITMQ_RETRY_DELAY_SECONDS", "5"))
	cronHourlyIntervalHours, _ := strconv.Atoi(getEnv("CRON_HOURLY_INTERVAL_HOURS", "1"))
	cronDailyIntervalHours, _ := strconv.Atoi(getEnv("CRON_DAILY_INTERVAL_HOURS", "24"))
	consumerLogIntervalSeconds, _ := strconv.Atoi(getEnv("CONSUMER_LOG_INTERVAL_SECONDS", "10"))

	return &Config{
		InfluxURL:    getEnv("INFLUX_URL", ""),
		InfluxToken:  getEnv("INFLUX_TOKEN", ""),
		InfluxOrg:    getEnv("INFLUX_ORG", ""),
		InfluxBucket: getEnv("INFLUX_BUCKET", ""),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		RabbitMQURL:           getEnv("RABBITMQ_URL", ""),
		RabbitMQQueueName:     getEnv("RABBITMQ_QUEUE_NAME", "electricity_queue"),
		RabbitMQRoutingKey:    getEnv("RABBITMQ_ROUTING_KEY", "electricity_metrics"),
		RabbitMQExchange:      getEnv("RABBITMQ_EXCHANGE", "amq.topic"),
		RabbitMQPrefetchCount: rabbitMQPrefetchCount,
		RabbitMQRetryDelay:    time.Duration(rabbitMQRetryDelaySeconds) * time.Second,

		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		PGHOST:     getEnv("PG_HOST", ""),
		PGPORT:     getEnv("PG_PORT", ""),
		PGUSER:     getEnv("PG_USER", ""),
		PGPASSWORD: getEnv("PG_PASS", ""),
		PGDATABASE: getEnv("PG_DTBS", ""),
		PGSSLMODE:  getEnv("PG_SSLMODE", "disable"),

		CORSAllowOrigins: parseStringSlice(getEnv("CORS_ALLOW_ORIGINS", "*")),
		CORSAllowMethods: parseStringSlice(getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS")),
		CORSAllowHeaders: parseStringSlice(getEnv("CORS_ALLOW_HEADERS", "Content-Type,Authorization")),

		CronHourlyInterval: time.Duration(cronHourlyIntervalHours) * time.Hour,
		CronDailyInterval:  time.Duration(cronDailyIntervalHours) * time.Hour,

		ConsumerLogInterval: time.Duration(consumerLogIntervalSeconds) * time.Second,

		SendgridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		SendgridFromEmail: getEnv("SENDGRID_FROM_EMAIL", ""),
		SendgridFromName: getEnv("SENDGRID_FROM_NAME", ""),
	}, nil
}

func parseStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
