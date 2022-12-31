package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type spotifyErrorData struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type spotifyError struct {
	Error spotifyErrorData `json:"error"`
}

func MakeRequest(method string, url string, body io.Reader, resultFormat interface{}) error {
	var errorResults = new(spotifyError)
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

	if err := json.NewDecoder(response.Body).Decode(errorResults); err != nil {
		return err
	}

	return fmt.Errorf("%s (%v)", errorResults.Error.Message, errorResults.Error.Status)
}
