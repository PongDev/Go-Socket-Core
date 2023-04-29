package services

import (
	"log"
	"net/http"

	services "github.com/PongDev/Go-Socket-Core/services/channel"
	"github.com/PongDev/Go-Socket-Core/types/dtos"
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

type WebsocketServiceInterface interface {
	HandleConnection(ctx *gin.Context)
}

type WebsocketService struct {
	channelService services.ChannelServiceInterface
}

func NewWebsocketService(channelService services.ChannelServiceInterface) WebsocketServiceInterface {
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
		var message dtos.SocketMessageDTO
		err := conn.ReadJSON(&message)

		if err != nil {
			s.channelService.DisconnectClient(conn)
			conn.Close()
			break
		}

		t := message.Type
		channelId := message.ChannelID

		switch t {
		case dtos.SocketMessageTypeJoin:
			s.channelService.JoinChannel(conn, channelId)
			conn.WriteJSON(dtos.SocketMessageDTO{
				Type: dtos.SocketMessageTypeACK,
			})
		case dtos.SocketMessageTypeLeave:
			s.channelService.LeaveChannel(conn, channelId)
			conn.WriteJSON(dtos.SocketMessageDTO{
				Type: dtos.SocketMessageTypeACK,
			})
		default:
			conn.WriteJSON(dtos.SocketMessageDTO{
				Type:    dtos.SocketMessageTypeError,
				Message: dtos.MessageInvalidMessageType,
			})
		}
	}
}
