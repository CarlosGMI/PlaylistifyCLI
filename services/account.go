package services

import (
	"errors"

	"github.com/CarlosGMI/Playlistify/utils"
)

func GetAccountInformation() (string, error) {
	if err := IsAuthenticated(); err == nil {
		return "", errors.New(utils.AlreadyLoggedInError)
	}

	return "", nil
}
