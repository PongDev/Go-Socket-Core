package types

type CreateChannelResponse struct {
	ChannelId string `json:"channelId"`
}

type CreateChannelWithIdResponse struct {
	Message string `json:"message"`
}

type DeleteChannelResponse struct {
	Message string `json:"message"`
}

type HandleMessageResponse struct {
	Message string `json:"message"`
}
