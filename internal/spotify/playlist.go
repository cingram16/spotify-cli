package spotify

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
)

func (c *Client) ListPlaylists(ctx context.Context) ([]*spotify.SimplePlaylist, error) {
	// Get the current user's playlists
	playlists, err := c.spotify.GetPlaylistsForUser(ctx, c.user.ID)
	if err != nil {
		return nil, fmt.Errorf("getting playlists: %w", err)
	}

	// Create a slice of Playlist structs
	var userPlaylists []*spotify.SimplePlaylist
	for _, playlist := range playlists.Playlists {
		userPlaylists = append(userPlaylists, &spotify.SimplePlaylist{
			ID:   playlist.ID,
			Name: playlist.Name,
		})
	}

	return userPlaylists, nil
}

func (c *Client) RotatePlaylist(ctx context.Context, playlist spotify.SimplePlaylist, seeds spotify.Seeds) error {
	err := c.fillPlaylistWithUnheardTracksBySeeds(ctx, playlist.ID, seeds)
	if err != nil {
		return fmt.Errorf("filling playlist: %w", err)
	}

	return nil
}

func (c *Client) fillPlaylistWithUnheardTracksBySeeds(ctx context.Context, playlistID spotify.ID, seeds spotify.Seeds) error {
	var selectedTracks []spotify.ID
	trackIDSet := make(map[string]struct{})
	artistIDSet := make(map[string]struct{})

	recentlyPlayedTracks, err := c.GetUserTopTracks(ctx)
	if err != nil {
		return fmt.Errorf("getting recently played tracks: %w", err)
	}

	topArtists, err := c.GetUserTopArtists(ctx)
	if err != nil {
		return fmt.Errorf("getting top artists: %w", err)
	}

	for _, track := range recentlyPlayedTracks {
		trackIDSet[track.ID.String()] = struct{}{}
	}
	for _, artist := range topArtists {
		artistIDSet[artist.ID.String()] = struct{}{}
	}

	for len(selectedTracks) < 50 {
		recommendations, err := c.spotify.GetRecommendations(ctx,
			seeds,
			&spotify.TrackAttributes{},
			spotify.Limit(50))
		if err != nil {
			return fmt.Errorf("getting recommendations: %w", err)
		}

		for _, track := range recommendations.Tracks {
			if _, trackExists := trackIDSet[track.ID.String()]; !trackExists {
				if _, artistExists := artistIDSet[track.Artists[0].ID.String()]; !artistExists {
					selectedTracks = append(selectedTracks, track.ID)
					if len(selectedTracks) >= 50 {
						break
					}
				}
			}
		}
	}

	err = c.spotify.ReplacePlaylistTracks(ctx, spotify.ID(playlistID), selectedTracks...)
	if err != nil {
		return fmt.Errorf("replacing playlist items: %w", err)
	}

	return nil
}
