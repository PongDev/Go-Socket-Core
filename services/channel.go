package services

import (
	"net/http"

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
	hub *utils.HubInterface
}

func NewChannelService(hub *utils.HubInterface) ChannelServiceInterface {
	return &ChannelService{
		hub: hub,
	}
}

func (s *ChannelService) CreateChannel(ctx *gin.Context) {

}

func (s *ChannelService) ConnectChannel(ctx *gin.Context) {

}

func (s *ChannelService) DeleteChannel(ctx *gin.Context) {

}

func (s *ChannelService) HandleMessage(ctx *gin.Context) {

}
