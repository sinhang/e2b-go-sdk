package e2b

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) ListTeams(ctx context.Context, query url.Values) ([]JSONMap, error) {
	var out []JSONMap
	err := c.doJSON(ctx, http.MethodGet, "/teams", query, nil, &out)
	return out, err
}

func (c *Client) GetTeamMetrics(ctx context.Context, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, "/teams/metrics", query, nil, &out)
	return out, err
}

func (c *Client) GetTeamMetricsMax(ctx context.Context, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, "/teams/metrics/max", query, nil, &out)
	return out, err
}
