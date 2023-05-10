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
	tea "github.com/charmbracelet/bubbletea"
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
type AuthMsg struct {
	ErrorType int
	Message   string
}
type NotAuthenticatedMsg AuthMsg
type AuthErrorMsg AuthMsg
type AuthorizedMsg string
type LoggedInMsg string
type LoggedInUserMsg struct{ Message string }

var pkceVerifier, pkceChallenge string
var authorization authorizationValues

func InitAuthentication() tea.Msg {
	token := viper.GetString("token")
	expiration := viper.GetInt64("token_expiration")

	if token == "" {
		return NotAuthenticatedMsg{utils.NotLoggedInCode, utils.NotLoggedInError}
	}

	if time.Now().Unix() > expiration {
		return NotAuthenticatedMsg{utils.ExpiredTokenCode, utils.ExpiredTokenError}
	}

	return AuthErrorMsg{utils.AlreadyLoggedInCode, utils.AlreadyLoggedInError}
}

func Authenticate() tea.Msg {
	initPKCECodeChallenge()
	state := generateRandomState()
	uri := buildAuthURI(pkceChallenge, state)

	if err := browser.OpenURL(uri); err != nil {
		return AuthErrorMsg{
			Message: err.Error(),
		}
	}

	authorization = listenForSpotifyAuthorization(state)

	if authorization.err != nil {
		return AuthErrorMsg{
			Message: authorization.err.Error(),
		}
	}

	return AuthorizedMsg("Authorized")
}

func Login() tea.Msg {
	token, err := requestSpotifyToken(authorization.code, pkceVerifier)

	if err != nil {
		return AuthErrorMsg{
			Message: err.Error(),
		}
	}

	storeTokenInformation(token)

	return LoggedInMsg("Authenticated")
}

func FetchAuthenticatedUser() tea.Msg {
	user, err := GetAccountInformation()

	if err != nil {
		return AuthErrorMsg{
			Message: err.Error(),
		}
	}

	return LoggedInUserMsg{fmt.Sprintf("Successfully logged in as %s (%s)", user.DisplayName, user.Email)}
}

func initPKCECodeChallenge() {
	verifier, _ := pkce.CreateCodeVerifier()
	pkceVerifier = verifier.Value
	pkceChallenge = verifier.CodeChallengeS256()
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

func RefreshAuthorization() tea.Msg {
	token, err := requestSpotifyRefreshToken()

	if err != nil {
		return AuthErrorMsg{
			Message: err.Error(),
		}
	}

	storeTokenInformation(token)

	return LoggedInMsg("Authenticated")
}

func requestSpotifyRefreshToken() (*token, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {viper.GetString("refresh_token")},
		"client_id":     {utils.ClientId},
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
	viper.Set("token_expiration", time.Now().Unix()+int64(token.ExpiresIn))
	viper.Set("refresh_token", token.RefreshToken)

	viper.WriteConfig()
}

func EmptyTokenInformation() {
	viper.Set("token", "")
	viper.Set("token_expiration", "")
	viper.Set("refresh_token", "")

	viper.WriteConfig()
}
