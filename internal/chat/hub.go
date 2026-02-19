package chat

import (
	"sync"
)

type Room struct {
	Clients map[*Client]bool
}

type Hub struct {
	Rooms map[string]*Room
	Mu    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) JoinRoom(roomName string, client *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	if _, exists := h.Rooms[roomName]; !exists {
		h.Rooms[roomName] = &Room{
			Clients: make(map[*Client]bool),
		}
	}

	h.Rooms[roomName].Clients[client] = true
}

func (h *Hub) LeaveRoom(client *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	room := h.Rooms[client.Room]
	if room == nil {
		return
	}

	delete(room.Clients, client)

	if len(room.Clients) == 0 {
		delete(h.Rooms, client.Room)
	}
}

func (h *Hub) Broadcast(roomName string, message []byte) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	room := h.Rooms[roomName]
	if room == nil {
		return
	}

	for client := range room.Clients {
		err := client.Conn.WriteMessage(1, message)
		if err != nil {
			client.Conn.Close()
			delete(room.Clients, client)
		}
	}
}
