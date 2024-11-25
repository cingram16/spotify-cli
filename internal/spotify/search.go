package spotify

import (
	"context"
	"github.com/zmb3/spotify/v2"
)

func (c *Client) Search(ctx context.Context, query string, t spotify.SearchType, limit int) (*spotify.SearchResult, error) {
	return c.spotify.Search(ctx, query, t, spotify.Limit(limit))
}

func (c *Client) GetFirstSearchArtist(ctx context.Context, query string) (*spotify.FullArtist, error) {
	result, err := c.Search(ctx, query, spotify.SearchTypeArtist, 1)
	if err != nil {
		return nil, err
	}

	if len(result.Artists.Artists) == 0 {
		return nil, nil
	}

	return &result.Artists.Artists[0], nil
}
