package main

import (
	"context"
	"log"

	"metertronik/internal/handler/amqp"
	"metertronik/internal/router"
	"metertronik/internal/service"
	"metertronik/pkg/config"
	"metertronik/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	influxRepo, cleanupInflux := database.SetupInfluxDB(cfg)
	defer cleanupInflux()

	RedisRealtimeRepo, cleanupRedis := database.SetupRedisRealtime(cfg)
	defer cleanupRedis()

	svc := service.NewIngestService(influxRepo, RedisRealtimeRepo)

	consumerCfg := &amqp.ConsumerConfig{
		QueueName:     cfg.RabbitMQQueueName,
		RoutingKey:    cfg.RabbitMQRoutingKey,
		Exchange:      cfg.RabbitMQExchange,
		PrefetchCount: cfg.RabbitMQPrefetchCount,
		RetryDelay:    cfg.RabbitMQRetryDelay,
		LogInterval:   cfg.ConsumerLogInterval,
	}
	consumer := amqp.NewConsumer(svc, consumerCfg)

	router.SetupWs(RedisRealtimeRepo, cfg.Port)

	ctx := context.Background()
	log.Printf("Consumer started, waiting for messages...")
	if err := consumer.StartConsuming(ctx, cfg.RabbitMQURL); err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
}
