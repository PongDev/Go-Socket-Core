package services

import (
	"log"
	"net/http"

	services "github.com/PongDev/Go-Socket-Core/services/channel"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketService struct {
	channelService services.ChannelServiceInterface
}

func NewWebsocketService(channelService services.ChannelServiceInterface) *WebsocketService {
	return &WebsocketService{
		channelService: channelService,
	}
}

func (s *WebsocketService) HandleConnection(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		var message map[string]interface{}
		err := conn.ReadJSON(&message)

		if err != nil {
			s.channelService.UnregisterClient(conn)
			conn.Close()
			break
		}

		t := message["type"]
		channelId := message["channelId"].(string)

		switch t {
		case "join":
			s.channelService.JoinChannel(conn, channelId)
			conn.WriteJSON(map[string]interface{}{
				"type": "ACK",
			})
		case "leave":
			s.channelService.LeaveChannel(conn, channelId)
		default:
			conn.WriteJSON(map[string]interface{}{
				"type":    "ERROR",
				"message": "Invalid message type",
			})
		}
	}
}
