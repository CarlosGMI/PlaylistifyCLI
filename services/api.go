package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type spotifyError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func MakeRequest(method string, url string, body io.Reader, resultFormat interface{}) error {
	var spotifyErrorFormat = new(spotifyError)
	client := &http.Client{}
	request, error := http.NewRequest(method, url, body)
	request.Close = true

	if error != nil {
		return error
	}

	request.Header.Add("Authorization", "Bearer "+viper.GetString("token"))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Encoding", "identity")

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		if err := json.NewDecoder(response.Body).Decode(resultFormat); err != nil {
			return err
		}

		return nil
	}

	if err := json.NewDecoder(response.Body).Decode(spotifyErrorFormat); err != nil {
		return err
	}

	return errors.New("An error has occurred: " + spotifyErrorFormat.Message + "(" + spotifyErrorFormat.Status + ")")
}
