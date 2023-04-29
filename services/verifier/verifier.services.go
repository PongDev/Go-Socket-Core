package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/PongDev/Go-Socket-Core/types/dtos"
)

func VerifyOperation(token string, channelId string, types dtos.SocketMessageType) bool {
	reqBody, err := json.Marshal(dtos.VerifierRequestDTO{
		Type:      types,
		ChannelID: channelId,
		Token:     token,
	})
	if err != nil {
		return false
	}
	req, err := http.NewRequest("POST", os.Getenv("VERIFIER_URL"), bytes.NewReader(reqBody))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return false
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}
	var verifier dtos.VerifierResponseDTO
	if err := json.Unmarshal(body, &verifier); err != nil {
		return false
	}
	return verifier.Valid
}
