package services

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lithammer/fuzzysearch/fuzzy"
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
	Id            string             `json:"id"`
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

type trackArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type trackInfo struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Artists []trackArtist `json:"artists"`
}

type track struct {
	Track trackInfo `json:"track"`
}

type playlistTracks struct {
	Tracks []track `json:"items"`
}

func GetPlaylists() error {
	var playlists []playlist
	var query = url.Values{
		"limit": {strconv.Itoa(utils.TracksLimit)},
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

func SearchInPlaylist(playlistId int, searchTerm string) error {
	var playlists []playlist
	var playlist = new(playlist)

	if err := viper.UnmarshalKey("playlists", &playlists); err != nil {
		return err
	}

	if len(playlists) > 0 && playlistId <= len(playlists) {
		playlist = &playlists[playlistId]
	} else {
		if err := getPlaylistWithOffset(strconv.Itoa(playlistId), playlist); err != nil {
			return err
		}
	}

	getTracksAndSearch(playlist, searchTerm)

	return nil
}

func getPlaylistWithOffset(id string, playlist *playlist) error {
	var playlistsResults = new(Playlists)
	var query = url.Values{
		"limit":  {"1"},
		"offset": {id},
	}
	var url = fmt.Sprintf("%s/me/playlists?%s", utils.SpotifyAPIBaseURL, query.Encode())

	if err := MakeRequest(http.MethodGet, url, nil, playlistsResults); err != nil {
		return err
	}

	if len(playlistsResults.Items) == 0 {
		return fmt.Errorf(utils.InexistentPlaylistError, id)
	}

	*playlist = playlistsResults.Items[0]

	return nil
}

func getTracksAndSearch(playlist *playlist, searchTerm string) error {
	var results []table.Row
	waitGroup := sync.WaitGroup{}
	numberOfRequests := math.Ceil(float64(playlist.Tracks.Total) / utils.TracksLimit)

	for i := 0; i < int(numberOfRequests); i++ {
		var tracksResults = new(playlistTracks)

		waitGroup.Add(1)

		go func(requestNumber int, playlistId string, tracks *playlistTracks, term string, results *[]table.Row) {
			fetchTracks(requestNumber, playlist.Id, tracks)
			executeSearch(tracks.Tracks, term, requestNumber, results)
			waitGroup.Done()
		}(i, playlist.Id, tracksResults, searchTerm, &results)
	}

	waitGroup.Wait()
	printSearchResults(results)

	return nil
}

func fetchTracks(requestNumber int, playlistId string, tracksResults *playlistTracks) error {
	query := url.Values{
		"limit":  {strconv.Itoa(utils.TracksLimit)},
		"offset": {strconv.Itoa(requestNumber * utils.TracksLimit)},
		"fields": {"items(track(name,id,artists(name,id)))"},
	}
	url := fmt.Sprintf("%s/playlists/%s/tracks?%s", utils.SpotifyAPIBaseURL, playlistId, query.Encode())

	if err := MakeRequest(http.MethodGet, url, nil, tracksResults); err != nil {
		return err
	}

	return nil
}

func executeSearch(tracks []track, term string, requestNumber int, results *[]table.Row) error {
	var offset = requestNumber * utils.TracksLimit

	for i, item := range tracks {
		var artists []string

		for _, artist := range item.Track.Artists {
			artists = append(artists, artist.Name)
		}

		formattedArtists := strings.Join(artists, ", ")
		trackNameIncludesTerm := fuzzy.Match(term, strings.ToLower(item.Track.Name))
		trackNameTermScore := CalculateJaroWinkler(term, strings.ToLower(item.Track.Name))
		artistsIncludesTerm := fuzzy.Match(term, strings.ToLower(formattedArtists))
		artistsTermScore := CalculateJaroWinkler(term, strings.ToLower(formattedArtists))

		if trackNameIncludesTerm || artistsIncludesTerm || trackNameTermScore > 0.8 || artistsTermScore > 0.8 {
			*results = append(*results, table.Row{strconv.Itoa(offset + i + 1), item.Track.Name, formattedArtists})
		}
	}

	return nil
}

func printSearchResults(results []table.Row) {
	resultsTable := table.NewWriter()

	resultsTable.SetOutputMirror(os.Stdout)
	resultsTable.AppendHeader(table.Row{"#", "Name", "Artists"})
	resultsTable.AppendRows(results)
	resultsTable.Render()
}
