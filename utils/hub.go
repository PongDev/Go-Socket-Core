package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/PongDev/Go-Socket-Core/types"
	"github.com/PongDev/Go-Socket-Core/types/dtos"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type HubInterface interface {
	CheckExistChannel(string) bool
	CreateChannel() string
	CloseChannel(string) error
	JoinChannel(string, *websocket.Conn) error
	LeaveChannel(string, *websocket.Conn) error
	DisconnectClient(*websocket.Conn) error
	SendMessageToChannel(string, []byte) error
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

func (h *Hub) CreateChannel() string {
	h.lock.Lock()
	defer h.lock.Unlock()

	var channelId string

	for {
		channelId = uuid.New().String()
		if _, ok := h.channels[channelId]; !ok {
			break
		}
	}

	h.channels[channelId] = make(map[*websocket.Conn]bool)

	return channelId
}

func (h *Hub) CloseChannel(channelId string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.channels[channelId]; !ok {
		return types.ChannelNotFoundError(channelId)
	}
	for conn := range h.channels[channelId] {
		err := conn.WriteJSON(dtos.SocketMessageDTO{
			Type:      dtos.SocketMessageTypeCloseChannel,
			ChannelID: channelId,
		})
		if err != nil {
			log.Println(err)
		}
		delete(h.clients, conn)
	}

	delete(h.channels, channelId)
	return nil
}

func (h *Hub) JoinChannel(channelId string, conn *websocket.Conn) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.channels[channelId]; !ok {
		return types.ChannelNotFoundError(channelId)
	}
	h.channels[channelId][conn] = true

	if h.clients[conn] == nil {
		h.clients[conn] = make(map[string]bool)
	}
	h.clients[conn][channelId] = true
	return nil
}

func (h *Hub) LeaveChannel(channelId string, conn *websocket.Conn) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.channels[channelId]; !ok {
		return types.ChannelNotFoundError(channelId)
	}
	if _, ok := h.clients[conn]; !ok {
		return types.ClientNotFoundError(fmt.Sprint(&conn))
	}
	delete(h.channels[channelId], conn)
	delete(h.clients[conn], channelId)
	return nil
}

func (h *Hub) DisconnectClient(conn *websocket.Conn) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.clients[conn]; !ok {
		return types.ClientNotFoundError(fmt.Sprint(&conn))
	}
	for channelId := range h.clients[conn] {
		delete(h.channels[channelId], conn)
	}
	delete(h.clients, conn)
	return nil
}

func (h *Hub) SendMessageToChannel(channelId string, message []byte) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.channels[channelId]; !ok {
		return types.ChannelNotFoundError(channelId)
	}
	for conn := range h.channels[channelId] {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (h *Hub) BroadcastMessage(message []byte) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
		}
	}
}
