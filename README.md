## Spotify CLI

## Installation

1. Install the application using `go install`:
    ```sh
    go install github.com/cingram16/spotify-cli@latest
    ```

## Spotify Developer App Setup

1. Create a developer app in the Spotify Developer Dashboard: [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications)

2. Set the environment variables with your Spotify Client ID and Secret:
    ```sh
    export SPOTIFY_ID=your_spotify_client_id
    export SPOTIFY_SECRET=your_spotify_secret
    ```

## Usage

Run the Spotify CLI:
```sh
spotify
```

### Example: Using the `playlist rotate` Command

The `playlist rotate` command allows you to select a playlist to update with tracks similar to the seed artists but avoids songs you commonly listen to.

```sh
spotify playlist rotate
```
