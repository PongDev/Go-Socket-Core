package utils

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type HubInterface interface {
	CheckExistChannel(string) bool
	JoinChannel(string, *websocket.Conn)
	LeaveChannel(string, *websocket.Conn)
	UnregisterClient(*websocket.Conn)
	RegisterNewChannel() string
	UnregisterChannel(string)
	SendMessageToChannel(string, []byte)
	BroadcastMessage([]byte)
}

type Hub struct {
	// Channel Id -> Connection[]
	channels map[string]map[*websocket.Conn]bool

	// Connection -> Channel Id
	clients map[*websocket.Conn]map[string]bool

	lock *sync.Mutex
}

func NewHub() HubInterface {
	return &Hub{
		channels: make(map[string]map[*websocket.Conn]bool),
		clients:  make(map[*websocket.Conn]map[string]bool),
		lock:     &sync.Mutex{},
	}
}

func (h *Hub) CheckExistChannel(channelId string) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.channels[channelId]
	return ok
}

func (h *Hub) JoinChannel(channelId string, conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.channels[channelId][conn] = true

	if h.clients[conn] == nil {
		h.clients[conn] = make(map[string]bool)
	}
	h.clients[conn][channelId] = true
}

func (h *Hub) LeaveChannel(channelId string, conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.channels[channelId], conn)
	delete(h.clients[conn], channelId)
}

func (h *Hub) RegisterNewChannel() string {
	h.lock.Lock()
	defer h.lock.Unlock()

	uuid := uuid.New()
	channelId := uuid.String()

	h.channels[channelId] = make(map[*websocket.Conn]bool)

	return channelId
}

func (h *Hub) SendMessageToChannel(channelId string, message []byte) {
	for conn := range h.channels[channelId] {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (h *Hub) BroadcastMessage(message []byte) {
	for conn := range h.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (h *Hub) UnregisterChannel(channelId string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for conn := range h.channels[channelId] {
		conn.WriteMessage(websocket.CloseMessage, []byte(fmt.Sprintf("Unregister channel: %s", channelId)))
		delete(h.clients, conn)
	}

	delete(h.channels, channelId)
}

func (h *Hub) UnregisterClient(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for channelId := range h.clients[conn] {
		delete(h.channels[channelId], conn)
	}

	delete(h.clients, conn)
}
