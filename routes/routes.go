package routes

import (
	"github.com/PongDev/Go-Socket-Core/services"
	"github.com/PongDev/Go-Socket-Core/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	hub := utils.NewHub()

	channelService := services.NewChannelService(hub)

	r.POST("/channel", channelService.CreateChannel)
	r.POST("/channel/:channelId", channelService.HandleMessage)
	r.GET("/channel/:channelId", channelService.ConnectChannel)
	r.DELETE("/channel/:channelId", channelService.DeleteChannel)
}
