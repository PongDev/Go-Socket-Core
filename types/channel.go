package types

type CreateChannelResponse struct {
	ChannelId string `json:"channelId"`
}

type DeleteChannelResponse struct {
	Message string `json:"message"`
}

type HandleMessageResponse struct {
	Message string `json:"message"`
}
