package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type Hub struct {
	Clients map[*Client]bool
	Lock    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	h.Clients[c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	delete(h.Clients, c)
	c.Conn.Close()
}

func (h *Hub) Broadcast(message []byte) {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	for c := range h.Clients {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("WebSocket error:", err)
			delete(h.Clients, c)
			c.Conn.Close()
		}
	}
}
