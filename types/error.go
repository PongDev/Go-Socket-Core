package types

import "fmt"

type JoinChannelError string

func (e JoinChannelError) Error() string {
	return fmt.Sprintf("Error Joining Channel ID: %s", string(e))
}

type ChannelNotFoundError string

func (e ChannelNotFoundError) Error() string {
	return fmt.Sprintf("Channel ID: %s Not Found", string(e))
}

type ClientNotFoundError string

func (e ClientNotFoundError) Error() string {
	return fmt.Sprintf("Client ID: %s Not Found", string(e))
}