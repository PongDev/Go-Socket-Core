package services

import (
	"log"
	"net/http"

	services "github.com/PongDev/Go-Socket-Core/services/channel"
	verifier "github.com/PongDev/Go-Socket-Core/services/verifier"
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
			err := s.channelService.DisconnectClient(conn)
			if err != nil {
				log.Println(err)
			}
			conn.Close()
			break
		}

		t := message.Type
		channelId := message.ChannelID

		switch t {
		case dtos.SocketMessageTypeJoin:
			if verifier.VerifyOperation(message.Token, channelId, dtos.SocketMessageTypeJoin) {
				err := s.channelService.JoinChannel(conn, channelId)
				if err != nil {
					log.Println(err)
				}
				err = conn.WriteJSON(dtos.SocketMessageDTO{
					Type: dtos.SocketMessageTypeACK,
				})
				if err != nil {
					log.Println(err)
				}
			} else {
				err := conn.WriteJSON(dtos.SocketMessageDTO{
					Type: dtos.SocketMessageTypeUnauthorized,
				})
				if err != nil {
					log.Println(err)
				}
			}
		case dtos.SocketMessageTypeLeave:
			err := s.channelService.LeaveChannel(conn, channelId)
			if err != nil {
				log.Println(err)
			}
			err = conn.WriteJSON(dtos.SocketMessageDTO{
				Type: dtos.SocketMessageTypeACK,
			})
			if err != nil {
				log.Println(err)
			}
		default:
			err := conn.WriteJSON(dtos.SocketMessageDTO{
				Type:    dtos.SocketMessageTypeError,
				Message: dtos.MessageInvalidMessageType,
			})
			if err != nil {
				log.Println(err)
			}
		}
	}
}
