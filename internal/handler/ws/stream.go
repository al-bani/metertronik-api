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
		return
	}
	defer conn.Close()

	log.Printf("\n\n[%s]\nWebSocket client connected (%s)", time.Now().Format("2006-01-02 15:04:05"), deviceID)

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
				continue
			}

			if data == nil {
				continue
			}

			dataJSON, err := json.Marshal(data)
			if err != nil {
				continue
			}
			currentHash := string(dataJSON)

			if currentHash == lastDataHash {
				continue
			}

			conn.SetWriteDeadline(utils.TimeNow().Time.Add(writeWait))
			if err := conn.WriteJSON(data); err != nil {
				return
			}

			lastDataHash = currentHash
			log.Printf("Send data to Websocket")

		case <-done:
			log.Printf("WebSocket closed")
			return
		case <-ctx.Done():
			log.Printf("WebSocket context cancelled")
			return
		}
	}
}
