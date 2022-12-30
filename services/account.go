package services

import (
	"net/http"

	"github.com/CarlosGMI/Playlistify/utils"
)

type UserAccount struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func GetAccountInformation() (*UserAccount, error) {
	if err := IsAuthenticated(); err != nil {
		return nil, err
	}

	var user = new(UserAccount)
	var url = utils.SpotifyAPIBaseURL + "/me"
	err := MakeRequest(http.MethodGet, url, nil, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
