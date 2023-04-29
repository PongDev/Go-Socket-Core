package routes

import (
	channelServices "github.com/PongDev/Go-Socket-Core/services/channel"
	websocketServices "github.com/PongDev/Go-Socket-Core/services/websocket"

	"github.com/PongDev/Go-Socket-Core/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	hub := utils.NewHub()

	channelService := channelServices.NewChannelService(hub)
	websocketService := websocketServices.NewWebsocketService(channelService)

	r.GET("/", websocketService.HandleConnection)
	r.POST("/channel", channelService.CreateChannel)
	r.POST("/channel/:channelId", channelService.CreateChannelWithId)
	r.POST("/channel/message/:channelId", channelService.HandleMessage)
	r.DELETE("/channel/:channelId", channelService.CloseChannel)
}
