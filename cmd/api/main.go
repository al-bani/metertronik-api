package main

import (
	"log"
	handler "metertronik/internal/handler/api"
	"metertronik/internal/middleware"
	httpRouter "metertronik/internal/router/http"
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

	api := service.NewApiService(postgresRepo, redisBatchRepo)
	apiHandler := handler.NewApiHandler(api)

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	httpRouter.SetupRoutes(router, apiHandler)
	router.Run(":" + cfg.Port)

	log.Printf("✅ API server started on port %s", cfg.Port)
}
