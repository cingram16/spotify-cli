package spotify

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
)

func (c *Client) GetUserTopTracks(ctx context.Context) ([]spotify.FullTrack, error) {
	var tracks []spotify.FullTrack

	trackPage, err := c.spotify.CurrentUsersTopTracks(ctx, spotify.Limit(50))
	if err != nil {
		return nil, fmt.Errorf("getting top played tracks: %w", err)
	}

	for _, track := range trackPage.Tracks {
		tracks = append(tracks, track)
	}

	return tracks, nil
}
