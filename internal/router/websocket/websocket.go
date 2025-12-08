package websocket

import (
	"net/http"

	"metertronik/internal/domain/repository"
	wsHandler "metertronik/internal/handler/ws"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(r *gin.Engine, RedisRepo repository.RedisRepo) {
	if RedisRepo == nil {
		return
	}

	wsStreamHandler := wsHandler.NewStreamHandler(RedisRepo)

	r.GET("/ws/electricity/:deviceID", func(c *gin.Context) {
		deviceID := c.Param("deviceID")
		if deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
			return
		}
		wsStreamHandler.HandleWebSocket(c.Writer, c.Request, deviceID)
	})
}
