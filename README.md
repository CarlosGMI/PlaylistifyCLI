# Playlistify

CLI application to search for a song or an artist in a Spotify playlist. This could work for you if you want to check if a specific song already exists in a playlist.

## Usage

### To login using your Spotify account

```bash
playlistify login
```

```bash
go run ./main.go login
```

### To logout from your Spotify account

```bash
playlistify logout
```

```bash
go run ./main.go logout
```

### To list all your playlists

```bash
playlistify list
```

```bash
go run ./main.go list
```

You could use this command to grab the playlist ID you want to perform the search in

### To search:

```bash
playlistify search
```


```bash
go run ./main.go search
```

#### Optional flags

- `--playlist`, `-p` | The playlist ID 
- `--term`, `-t` | The term you want to search in the playlist

```bash
playlistify search -p 10 -t "Term"
```

```bash
go run ./main.go search -p 10 -t "Term"
```
