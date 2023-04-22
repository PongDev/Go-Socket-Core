package services

import (
	"log"
	"net/http"

	"github.com/PongDev/Go-Socket-Core/dtos"
	"github.com/PongDev/Go-Socket-Core/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChannelServiceInterface interface {
	CreateChannel(ctx *gin.Context)
	DeleteChannel(ctx *gin.Context)
	HandleMessage(ctx *gin.Context)
	UnregisterClient(conn *websocket.Conn)
	JoinChannel(conn *websocket.Conn, channelId string)
	LeaveChannel(conn *websocket.Conn, channelId string)
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

func (s *ChannelService) JoinChannel(conn *websocket.Conn, channelId string) {
	s.hub.JoinChannel(channelId, conn)
}

func (s *ChannelService) LeaveChannel(conn *websocket.Conn, channelId string) {
	s.hub.LeaveChannel(channelId, conn)
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

func (s *ChannelService) UnregisterClient(conn *websocket.Conn) {
	s.hub.UnregisterClient(conn)
}
