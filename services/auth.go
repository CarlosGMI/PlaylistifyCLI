package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/CarlosGMI/Playlistify/utils"
	pkce "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

type authorizationValues struct {
	code string
	err  error
}

type token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func Authenticate() error {
	if err := IsAuthenticated(); err == nil {
		user, err := GetAccountInformation()

		if err != nil {
			return err
		}

		return fmt.Errorf(utils.AlreadyLoggedInError, user.DisplayName, user.Email)
	}

	if err := login(); err != nil {
		return err
	}

	return nil
}

func IsAuthenticated() error {
	token := viper.GetString("token")

	if token == "" {
		return errors.New(utils.NotLoggedInError)
	}

	// Token exists -> check if token has expired
	// If expired, refresh token here. If not, return nil

	return nil
}

func login() error {
	pkceVerifier, pkceChallenge := initPKCECodeChallenge()
	state := generateRandomState()
	uri := buildAuthURI(pkceChallenge, state)

	if err := browser.OpenURL(uri); err != nil {
		return err
	}

	authorizationValues := listenForSpotifyAuthorization(state)

	if authorizationValues.err != nil {
		return authorizationValues.err
	}

	fmt.Println("Authenticating...")

	token, err := requestSpotifyToken(authorizationValues.code, pkceVerifier)

	if err != nil {
		return err
	}

	storeTokenInformation(token)
	user, err := GetAccountInformation()

	if err != nil {
		return err
	}

	fmt.Printf("Successfully logged in as %s (%s)\n", user.DisplayName, user.Email)

	return nil
}

func initPKCECodeChallenge() (string, string) {
	pkceVerifier, _ := pkce.CreateCodeVerifier()
	pkceChallenge := pkceVerifier.CodeChallengeS256()

	return pkceVerifier.Value, pkceChallenge
}

func generateRandomState() string {
	rand.Seed(time.Now().Unix())

	letters := []rune(utils.LetterRunes)
	state := make([]rune, 15)

	for i := range state {
		state[i] = letters[rand.Intn(len(utils.LetterRunes))]
	}

	return string(state)
}

func buildAuthURI(pkceChallenge string, state string) string {
	queryParams := url.Values{
		"client_id":             {utils.ClientId},
		"response_type":         {"code"},
		"redirect_uri":          {utils.AuthorizationCallbackURL},
		"state":                 {state},
		"scope":                 {utils.PlaylistifyScopes},
		"code_challenge_method": {"S256"},
		"code_challenge":        {pkceChallenge},
		"show_dialog":           {"false"},
	}

	return utils.SpotifyAccountBaseURL + "/authorize?" + queryParams.Encode()
}

func listenForSpotifyAuthorization(state string) authorizationValues {
	var values = authorizationValues{
		code: "",
		err:  nil,
	}
	server := &http.Server{Addr: utils.AuthorizationPort}

	http.HandleFunc("/callback", func(writer http.ResponseWriter, req *http.Request) {
		code, err := getSpotifyAuthorization(writer, req, state)
		values.code = code
		values.err = err

		go func() {
			server.Shutdown(context.Background())
		}()
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return values
	}

	return values
}

func getSpotifyAuthorization(writer http.ResponseWriter, req *http.Request, state string) (string, error) {
	url, _ := url.Parse(req.URL.String())
	queryParams := url.Query()
	code := queryParams.Get("code")
	spotifyState := queryParams.Get("state")

	if queryParams.Has("error") || spotifyState != state {
		fmt.Fprintf(writer, "An error occured: "+queryParams.Get("error"))

		return "", errors.New(utils.NotAuthorizedError)
	} else {
		fmt.Fprintf(writer, "Authorization successful!")

		return code, nil
	}
}

func requestSpotifyToken(code string, pkceVerifier string) (*token, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {utils.AuthorizationCallbackURL},
		"client_id":     {utils.ClientId},
		"code_verifier": {pkceVerifier},
	}

	response, err := http.PostForm(utils.SpotifyAccountBaseURL+"/api/token", data)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	token := new(token)
	err = json.NewDecoder(response.Body).Decode(token)

	return token, err
}

func storeTokenInformation(token *token) {
	viper.Set("token", token.AccessToken)
	viper.Set("token_expiration", token.ExpiresIn)
	viper.Set("refresh_token", token.RefreshToken)

	viper.WriteConfig()
}
