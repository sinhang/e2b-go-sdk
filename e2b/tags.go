package e2b

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) AssignTags(ctx context.Context, templateID string, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/templates/%s/tags", templateID), nil, body, &out)
	return out, err
}

func (c *Client) DeleteTags(ctx context.Context, templateID string, body JSONMap) error {
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/templates/%s/tags", templateID), nil, body, nil)
}

func (c *Client) ListTemplateTags(ctx context.Context, templateID string, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/templates/%s/tags", templateID), query, nil, &out)
	return out, err
}
