package utils

import "github.com/charmbracelet/lipgloss"

const (
	// Errors
	NotLoggedInError        = `you are not logged in, please run "playlistify login"`
	AlreadyLoggedInError    = "you are already logged in as %s (%s)"
	NotAuthorizedError      = "you are not authorized"
	ExpiredTokenError       = "the authentication token has expired"
	InexistentPlaylistError = "playlist with ID of %s doesn't exist"
	NotLoggedInCode         = 0
	ExpiredTokenCode        = 1
	AlreadyLoggedInCode     = 2
	// General
	ClientId                      = "c4ab33f93b55422bb1cf39494023da7d"
	SpotifyAccountBaseURL         = "https://accounts.spotify.com"
	SpotifyAPIBaseURL             = "https://api.spotify.com/v1"
	AuthorizationBaseURL          = "http://localhost"
	AuthorizationPort             = ":1024"
	AuthorizationCallbackEndpoint = "/callback"
	AuthorizationCallbackURL      = AuthorizationBaseURL + AuthorizationPort + AuthorizationCallbackEndpoint
	LetterRunes                   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	PlaylistifyScopes             = "playlist-read-private playlist-read-collaborative user-read-email user-read-private"
	TracksLimit                   = 50
	// TUI Colors
	ColorSpotifyGreen = "#1DB954"
	ColorSpotifyRed   = "#FF5263"
	// TUI States
	LoadingState = "loading"
	ErrorState   = "error"
	SuccessState = "success"
)

var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSpotifyRed)).Render
