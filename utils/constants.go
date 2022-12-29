package utils

const (
	// Errors
	NotLoggedInError     = `you are not logged in, please run "playlistify login"`
	AlreadyLoggedInError = "you are already logged in"
	NotAuthorizedError   = "you are not authorized"
	// General
	ClientId                      = "c4ab33f93b55422bb1cf39494023da7d"
	SpotifyAccountBaseURL         = "https://accounts.spotify.com"
	SpotifyAPIBaseURL             = "https://api.spotify.com/v1"
	AuthorizationBaseURL          = "http://localhost"
	AuthorizationPort             = ":1024"
	AuthorizationCallbackEndpoint = "/callback"
	AuthorizationCallbackURL      = AuthorizationBaseURL + AuthorizationPort + AuthorizationCallbackEndpoint
	LetterRunes                   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	PlaylistifyScopes             = "playlist-read-private playlist-read-collaborative"
)
