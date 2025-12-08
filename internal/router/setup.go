package router

import (
	"log"

	"metertronik/internal/domain/repository"
	"metertronik/internal/router/websocket"

	"github.com/gin-gonic/gin"
)

func SetupWs(RedisRepo repository.RedisRepo, port string) {
	if RedisRepo == nil {
		log.Println("Warning: RedisRepo is nil, WebSocket server will not start")
		return
	}

	r := gin.Default()
	websocket.WebSocketRoutes(r, RedisRepo)

	go func() {
		log.Printf("🚀 WebSocket server starting on port %s", port)
		log.Printf("📡 WebSocket endpoint: ws://localhost:%s/ws/electricity/:deviceID", port)
		if err := r.Run(":" + port); err != nil {
			log.Printf("Failed to start WebSocket server: %v", err)
		}
	}()
}
