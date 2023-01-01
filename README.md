# Playlistify

CLI application to search for a song or an artist in a Spotify playlist. This could work for you if you want to check if a specific song already exists in a playlist.

## Usage

To login using your Spotify account

```bash
playlistify login
```

```bash
go run ./main.go login
```

To list all your playlists

```bash
playlistify list
```

```bash
go run ./main.go list
```

It's important to list the playlists so you can grab the ID of the playlist you want to search something in. To search:

```bash
playlistify search -p 10 -t "Term"
```

```bash
go run ./main.go search -p 10 -t "Term"
```

--playlist, -p | (required) The playlist ID --term, -t | (required) The term you want to search in the playlist

Additionally, you can logout of your account

```bash
playlistify logout
```

```bash
go run ./main.go logout
```
