package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
)

type playlistOwner struct {
	Id string `json:"id"`
}

type playlistTracksInfo struct {
	Total int    `json:"total"`
	Url   string `json:"href"`
}

type playlist struct {
	Collaborative bool               `json:"collaborative"`
	Name          string             `json:"name"`
	Type          string             `json:"type"`
	Owner         playlistOwner      `json:"owner"`
	Tracks        playlistTracksInfo `json:"tracks"`
}

type Playlists struct {
	Next  interface{} `json:"next"`
	Items []playlist  `json:"items"`
}

func GetPlaylists() error {

	var playlists []playlist
	var query = url.Values{
		"limit": {"50"},
	}
	var url = fmt.Sprintf("%s/me/playlists?%s", utils.SpotifyAPIBaseURL, query.Encode())

	if err := fetchPlaylists(&playlists, url); err != nil {
		return err
	}

	storePlaylists(&playlists)

	return nil
}

func fetchPlaylists(playlists *[]playlist, url string) error {
	var playlistsResults = new(Playlists)

	if err := MakeRequest(http.MethodGet, url, nil, playlistsResults); err != nil {
		return err
	}

	*playlists = append(*playlists, playlistsResults.Items...)

	if playlistsResults.Next != nil {
		if err := fetchPlaylists(playlists, playlistsResults.Next.(string)); err != nil {
			return nil
		}
	}

	return nil
}

func storePlaylists(playlists *[]playlist) {
	viper.Set("playlists", playlists)
	viper.WriteConfig()
}

func PrintPlaylists() error {
	var playlists []playlist
	userId := viper.GetString("user_id")

	if err := viper.UnmarshalKey("playlists", &playlists); err != nil {
		return err
	}

	playlistsTable := table.NewWriter()
	var rows []table.Row

	playlistsTable.SetOutputMirror(os.Stdout)
	playlistsTable.AppendHeader(table.Row{"Playlist ID", "Name", "Total Tracks"})

	for index, playlist := range playlists {
		if playlist.Owner.Id == userId || playlist.Collaborative {
			rows = append(rows, table.Row{index, playlist.Name, playlist.Tracks.Total})
		}
	}

	playlistsTable.AppendRows(rows)
	playlistsTable.Render()

	return nil
}
