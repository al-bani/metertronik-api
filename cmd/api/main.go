package main

import (
	"log"
	handler "metertronik/internal/handler/api"
	"metertronik/internal/middleware"
	httpRouter "metertronik/internal/router/http"
	wsRouter "metertronik/internal/router/websocket"
	"metertronik/internal/service"
	"metertronik/pkg/config"
	"metertronik/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	postgresRepo, cleanupPostgres := database.SetupPostgres(cfg)
	defer cleanupPostgres()

	redisBatchRepo, cleanupRedisBatch := database.SetupRedisRealtimeBatch(cfg)
	defer cleanupRedisBatch()

	redisRealtimeRepo, cleanupRedisRealtime := database.SetupRedisRealtime(cfg)
	defer cleanupRedisRealtime()

	api := service.NewApiService(postgresRepo, redisBatchRepo)
	apiHandler := handler.NewApiHandler(api)

	gin.SetMode(cfg.GinMode)
	router := gin.Default()

	router.Use(middleware.CORSMiddleware(cfg))

	httpRouter.SetupRoutes(router, apiHandler)

	wsRouter.WebSocketRoutes(router, redisRealtimeRepo)

	log.Printf("API server started on port %s", cfg.Port)
	log.Printf("HTTP API endpoint: http://localhost:%s/v1/api", cfg.Port)
	log.Printf("WebSocket endpoint: ws://localhost:%s/v1/ws/electricity/:deviceID", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
