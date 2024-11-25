package spotify

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
)

func (c *Client) GetUserTopArtists(ctx context.Context) ([]spotify.FullArtist, error) {
	var artists []spotify.FullArtist

	for i := 0; i < 5; i++ {
		artistPage, err := c.spotify.CurrentUsersTopArtists(ctx, spotify.Limit(50), spotify.Offset(i*50), spotify.Timerange(spotify.MediumTermRange))
		if err != nil {
			return nil, fmt.Errorf("getting top artists: %w", err)
		}

		for _, artist := range artistPage.Artists {
			artists = append(artists, artist)
		}
	}

	return artists, nil
}
