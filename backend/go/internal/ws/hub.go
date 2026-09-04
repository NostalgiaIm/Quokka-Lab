package ws

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Development deployments are proxied through Nginx; origin enforcement
		// should be tightened when production hostnames are known.
		return true
	},
}

// Hub maintains lightweight collaboration rooms for MIDI note events.
type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan subscription
	unregister chan subscription
	broadcast  chan roomMessage
	logger     *slog.Logger
	mu         sync.RWMutex
}

type subscription struct {
	roomID string
	client *Client
}

type roomMessage struct {
	roomID  string
	payload []byte
	sender  *Client
}

// NewHub creates a WebSocket hub with per-room fan-out.
func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		broadcast:  make(chan roomMessage, 256),
		logger:     logger,
	}
}

// Run processes hub events. It should be launched once at gateway startup.
func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			h.mu.Lock()
			if h.rooms[sub.roomID] == nil {
				h.rooms[sub.roomID] = make(map[*Client]bool)
			}
			h.rooms[sub.roomID][sub.client] = true
			h.mu.Unlock()
			h.logger.Info("websocket joined", "room", sub.roomID)

		case sub := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[sub.roomID]; ok {
				if _, exists := clients[sub.client]; exists {
					delete(clients, sub.client)
					close(sub.client.send)
				}
				if len(clients) == 0 {
					delete(h.rooms, sub.roomID)
				}
			}
			h.mu.Unlock()
			h.logger.Info("websocket left", "room", sub.roomID)

		case msg := <-h.broadcast:
			h.mu.Lock()
			for client := range h.rooms[msg.roomID] {
				if client == msg.sender {
					continue
				}
				select {
				case client.send <- msg.payload:
				default:
					close(client.send)
					delete(h.rooms[msg.roomID], client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// ServeRoom upgrades a browser request and attaches it to a collaboration room.
func (h *Hub) ServeRoom(w http.ResponseWriter, r *http.Request, roomID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", "error", err)
		return
	}

	client := &Client{hub: h, roomID: roomID, conn: conn, send: make(chan []byte, 256)}
	h.register <- subscription{roomID: roomID, client: client}

	go client.writePump()
	go client.readPump()
}
