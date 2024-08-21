package ws

import (
	wsmodels "chat-application/internal/websocket/models"
	"fmt"
)

type Hub struct {
	Rooms      map[string]*wsmodels.Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *wsmodels.Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*wsmodels.Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *wsmodels.Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomID]; ok {
				r := h.Rooms[client.RoomID]

				if _, ok = r.Clients[client.ID]; !ok {
					r.Clients[client.ID] = client
				}

			}
		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomID]; ok {
				if _, ok = h.Rooms[client.RoomID].Clients[client.ID]; ok {
					if len(h.Rooms[client.RoomID].Clients) > 0 {
						h.Broadcast <- &wsmodels.Message{
							Content:  fmt.Sprintf("%s left the chat", client.Username),
							RoomID:   client.RoomID,
							Username: client.Username,
						}
					}

					delete(h.Rooms[client.RoomID].Clients, client.ID)
					close(client.Message)
				}
			}
		case msg := <-h.Broadcast:
			if _, ok := h.Rooms[msg.RoomID]; ok {

				for _, cl := range h.Rooms[msg.RoomID].Clients {
					cl.Message <- msg
				}
			}
		}
	}
}
