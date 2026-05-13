package e2b

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetVolumes(ctx context.Context, query url.Values) ([]JSONMap, error) {
	var out []JSONMap
	err := c.doJSON(ctx, http.MethodGet, "/volumes", query, nil, &out)
	return out, err
}

func (c *Client) PostVolumes(ctx context.Context, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPost, "/volumes", nil, body, &out)
	return out, err
}

func (c *Client) GetVolume(ctx context.Context, volumeID string) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/volumes/%s", volumeID), nil, nil, &out)
	return out, err
}

func (c *Client) DeleteVolumes(ctx context.Context, volumeID string) error {
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/volumes/%s", volumeID), nil, nil, nil)
}
