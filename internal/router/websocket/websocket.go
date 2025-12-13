package websocket

import (
	"net/http"

	"metertronik/internal/domain/repository"
	wsHandler "metertronik/internal/handler/ws"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(r *gin.Engine, RedisRealtimeRepo repository.RedisRealtimeRepo) {
	if RedisRealtimeRepo == nil {
		return
	}

	wsStreamHandler := wsHandler.NewStreamHandler(RedisRealtimeRepo)

	r.GET("/ws/electricity/:deviceID", func(c *gin.Context) {
		deviceID := c.Param("deviceID")
		if deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
			return
		}
		wsStreamHandler.HandleWebSocket(c.Writer, c.Request, deviceID)
	})
}
