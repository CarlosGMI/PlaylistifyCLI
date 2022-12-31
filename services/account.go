package services

import (
	"net/http"

	"github.com/CarlosGMI/Playlistify/utils"
)

type UserAccount struct {
	Id          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func GetAccountInformation() (*UserAccount, error) {
	var user = new(UserAccount)
	var url = utils.SpotifyAPIBaseURL + "/me"
	err := MakeRequest(http.MethodGet, url, nil, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
