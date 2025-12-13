package ws

import (
	"context"
	"encoding/json"
	"log"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/utils"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	pollPeriod = 1 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type StreamHandler struct {
	RedisRealtimeRepo repository.RedisRealtimeRepo
}

func NewStreamHandler(RedisRealtimeRepo repository.RedisRealtimeRepo) *StreamHandler {
	return &StreamHandler{
		RedisRealtimeRepo: RedisRealtimeRepo,
	}
}

func (h *StreamHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request, deviceID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket client connected for device: %s", deviceID)

	conn.SetReadDeadline(utils.TimeNow().Time.Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(utils.TimeNow().Time.Add(pongWait))
		return nil
	})

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	dataTicker := time.NewTicker(pollPeriod)
	defer dataTicker.Stop()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case <-pingTicker.C:
				conn.SetWriteDeadline(utils.TimeNow().Time.Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Failed to send ping: %v", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	var lastDataHash string

	for {
		select {
		case <-dataTicker.C:
			data, err := h.RedisRealtimeRepo.GetLatestElectricity(ctx, deviceID)
			if err != nil {
				log.Printf("Error getting latest electricity data: %v", err)
				continue
			}

			if data == nil {
				continue
			}

			dataJSON, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling data: %v", err)
				continue
			}
			currentHash := string(dataJSON)

			if currentHash == lastDataHash {
				continue
			}

			conn.SetWriteDeadline(utils.TimeNow().Time.Add(writeWait))
			if err := conn.WriteJSON(data); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

			lastDataHash = currentHash
			log.Printf("Sent data update for device %s", deviceID)

		case <-done:
			log.Printf("WebSocket connection closed for device: %s", deviceID)
			return
		case <-ctx.Done():
			log.Printf("Context cancelled for device: %s", deviceID)
			return
		}
	}
}
