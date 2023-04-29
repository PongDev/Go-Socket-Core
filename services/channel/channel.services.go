package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/PongDev/Go-Socket-Core/types"
	"github.com/PongDev/Go-Socket-Core/types/dtos"
	"github.com/PongDev/Go-Socket-Core/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChannelServiceInterface interface {
	CreateChannel(ctx *gin.Context)
	CreateChannelWithId(ctx *gin.Context)
	CloseChannel(ctx *gin.Context)
	HandleMessage(ctx *gin.Context)
	JoinChannel(conn *websocket.Conn, channelId string) error
	LeaveChannel(conn *websocket.Conn, channelId string) error
	DisconnectClient(conn *websocket.Conn) error
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
	channelId := s.hub.CreateChannel()

	ctx.JSON(http.StatusOK, types.CreateChannelResponse{
		ChannelId: channelId,
	})
}

func (s *ChannelService) CreateChannelWithId(ctx *gin.Context) {
	channelId, ok := ctx.Params.Get("channelId")

	if !ok {
		ctx.JSON(http.StatusBadRequest, types.CreateChannelWithIdResponse{
			Message: dtos.MessageChannelIDRequired,
		})
		return
	}
	if err := s.hub.CreateChannelWithId(channelId); err != nil {
		ctx.JSON(http.StatusConflict, types.CreateChannelWithIdResponse{
			Message: dtos.MessageChannelExists,
		})
	}
	ctx.JSON(http.StatusOK, types.CreateChannelWithIdResponse{
		Message: dtos.MessageChannelCreated,
	})
}

func (s *ChannelService) CloseChannel(ctx *gin.Context) {
	channelId, ok := ctx.Params.Get("channelId")

	if !ok {
		ctx.JSON(http.StatusBadRequest, types.DeleteChannelResponse{
			Message: dtos.MessageChannelIDRequired,
		})
		return
	}
	if err := s.hub.CloseChannel(channelId); err != nil {
		var ChannelNotFoundError *types.ChannelNotFoundError
		if errors.As(err, &ChannelNotFoundError) {
			ctx.JSON(http.StatusInternalServerError, types.DeleteChannelResponse{
				Message: ChannelNotFoundError.Error(),
			})
		} else {
			ctx.JSON(http.StatusNotFound, types.DeleteChannelResponse{
				Message: dtos.MessageUnexpectedError,
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, types.DeleteChannelResponse{
		Message: dtos.MessageChannelDeleted,
	})
}

func (s *ChannelService) HandleMessage(ctx *gin.Context) {
	channelId, ok := ctx.Params.Get("channelId")

	if !ok {
		ctx.JSON(http.StatusBadRequest, types.HandleMessageResponse{
			Message: dtos.MessageChannelIDRequired,
		})
		return
	}

	var msg dtos.MessageDTO

	if err := ctx.BindJSON(&msg); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, types.HandleMessageResponse{
			Message: dtos.MessageInvalidRequest,
		})
		return
	}
	if err := s.hub.SendMessageToChannel(channelId, msg.Content); err != nil {
		var ChannelNotFoundError *types.ChannelNotFoundError
		if errors.As(err, &ChannelNotFoundError) {
			ctx.JSON(http.StatusNotFound, types.HandleMessageResponse{
				Message: ChannelNotFoundError.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, types.HandleMessageResponse{
				Message: dtos.MessageUnexpectedError,
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, types.HandleMessageResponse{
		Message: dtos.MessageSent,
	})
}

func (s *ChannelService) JoinChannel(conn *websocket.Conn, channelId string) error {
	if err := s.hub.JoinChannel(channelId, conn); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *ChannelService) LeaveChannel(conn *websocket.Conn, channelId string) error {
	err := s.hub.LeaveChannel(channelId, conn)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *ChannelService) DisconnectClient(conn *websocket.Conn) error {
	err := s.hub.DisconnectClient(conn)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
