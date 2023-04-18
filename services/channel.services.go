package services

import (
	"log"
	"net/http"

	"github.com/PongDev/Go-Socket-Core/dtos"
	"github.com/PongDev/Go-Socket-Core/utils"
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

type ChannelServiceInterface interface {
	CreateChannel(ctx *gin.Context)
	ConnectChannel(ctx *gin.Context)
	DeleteChannel(ctx *gin.Context)
	HandleMessage(ctx *gin.Context)
}

type ChannelService struct {
	hub utils.HubInterface
}

func NewChannelService(hub utils.HubInterface) ChannelServiceInterface {
	return &ChannelService{
		hub: hub,
	}
}

func (s *ChannelService) CreateChannel(ctx *gin.Context) {
	channelId := s.hub.RegisterNewChannel()

	ctx.JSON(http.StatusOK, gin.H{
		"channelId": channelId,
	})
}

func (s *ChannelService) ConnectChannel(ctx *gin.Context) {
	channelId := ctx.Param("channelId")

	if !s.hub.CheckExistChannel(channelId) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Channel not found",
		})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	s.hub.JoinChannel(channelId, conn)

	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			s.hub.UnregisterClient(conn)
			break
		}
	}
}

func (s *ChannelService) DeleteChannel(ctx *gin.Context) {
	channelId, _ := ctx.Params.Get("channelId")
	if !s.hub.CheckExistChannel(channelId) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Channel not found",
		})
		return
	}

	s.hub.UnregisterChannel(channelId)
}

func (s *ChannelService) HandleMessage(ctx *gin.Context) {
	channelId, _ := ctx.Params.Get("channelId")
	if !s.hub.CheckExistChannel(channelId) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Channel not found",
		})
		return
	}

	var msg dtos.Message

	err := ctx.BindJSON(&msg)
	if err != nil {
		log.Println(err)
		return
	}

	s.hub.SendMessageToChannel(channelId, []byte(msg.Content))
}
