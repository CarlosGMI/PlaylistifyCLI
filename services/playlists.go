package services

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	textTable "github.com/jedib0t/go-pretty/v6/table"

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

type PlaylistsErrorMsg struct {
	Message string
}

type PlaylistsMsg string
type SearchResultsMsg struct {
	Results     []table.Row
	TextResults []textTable.Row
}

func GetPlaylists() tea.Msg {
	var playlists []playlist
	var query = url.Values{
		"limit": {strconv.Itoa(utils.TracksLimit)},
	}
	var url = fmt.Sprintf("%s/me/playlists?%s", utils.SpotifyAPIBaseURL, query.Encode())

	if err := fetchPlaylists(&playlists, url); err != nil {
		return PlaylistsErrorMsg{err.Error()}
	}

	storePlaylists(&playlists)

	return PlaylistsMsg("")
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

func PrintPlaylists() ([]table.Row, []textTable.Row, error) {
	var playlists []playlist
	var rows []table.Row
	var textRows []textTable.Row
	userId := viper.GetString("user_id")

	if err := viper.UnmarshalKey("playlists", &playlists); err != nil {
		return rows, textRows, err
	}

	for index, playlist := range playlists {
		if playlist.Owner.Id == userId || playlist.Collaborative {
			rows = append(rows, table.Row{strconv.Itoa(index), playlist.Name, strconv.Itoa(playlist.Tracks.Total)})
			textRows = append(textRows, textTable.Row{strconv.Itoa(index), playlist.Name, strconv.Itoa(playlist.Tracks.Total)})
		}
	}

	return rows, textRows, nil
}

func SearchInPlaylist(playlistId string, searchTerm string) tea.Msg {
	var playlists []playlist
	var playlist = new(playlist)
	formattedId, err := strconv.Atoi(playlistId)

	if err != nil {
		return PlaylistsErrorMsg{err.Error()}
	}

	if err := viper.UnmarshalKey("playlists", &playlists); err != nil {
		return PlaylistsErrorMsg{err.Error()}
	}

	if len(playlists) > 0 && formattedId <= len(playlists) {
		playlist = &playlists[formattedId]
	} else {
		if err := getPlaylistWithOffset(strconv.Itoa(formattedId), playlist); err != nil {
			return PlaylistsErrorMsg{err.Error()}
		}
	}

	results, textResults, _ := getTracksAndSearch(playlist, strings.ToLower(searchTerm))

	return SearchResultsMsg{results, textResults}
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

func getTracksAndSearch(playlist *playlist, searchTerm string) ([]table.Row, []textTable.Row, error) {
	var results []table.Row
	var textResults []textTable.Row
	waitGroup := sync.WaitGroup{}
	numberOfRequests := math.Ceil(float64(playlist.Tracks.Total) / utils.TracksLimit)

	for i := 0; i < int(numberOfRequests); i++ {
		var tracksResults = new(playlistTracks)

		waitGroup.Add(1)

		go func(requestNumber int,
			playlistId string,
			tracks *playlistTracks,
			term string,
			results *[]table.Row,
			textResults *[]textTable.Row,
		) {
			fetchTracks(requestNumber, playlist.Id, tracks)
			executeSearch(tracks.Tracks, term, requestNumber, results, textResults)
			waitGroup.Done()
		}(i, playlist.Id, tracksResults, searchTerm, &results, &textResults)
	}

	waitGroup.Wait()

	return results, textResults, nil
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

func executeSearch(tracks []track, term string, requestNumber int, results *[]table.Row, textResults *[]textTable.Row) error {
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
			*textResults = append(*textResults, textTable.Row{strconv.Itoa(offset + i + 1), item.Track.Name, formattedArtists})
		}
	}

	return nil
}
