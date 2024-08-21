package wsmodels

import (
	ws "chat-application/internal/websocket"
)

type Room struct {
	ID      string                `json:"id"`
	Name    string                `json:"name"`
	Clients map[string]*ws.Client `json:"clients"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}
