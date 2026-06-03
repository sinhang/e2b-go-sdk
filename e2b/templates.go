package e2b

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) CreateTemplateV3(ctx context.Context, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	// Cube compat: /v3/templates is pure-404 on Cube, use legacy route directly.
	if c.compatMode {
		return out, c.doJSON(ctx, http.MethodPost, "/templates", nil, body, &out)
	}

	// E2B native: try /v3 → /v2 → legacy fallback chain.
	err := c.doJSON(ctx, http.MethodPost, "/v3/templates", nil, body, &out)
	if err == nil {
		return out, nil
	}
	if apiErr, ok := err.(*APIResponseError); !ok || apiErr.StatusCode != http.StatusNotFound {
		return out, err
	}

	err = c.doJSON(ctx, http.MethodPost, "/v2/templates", nil, body, &out)
	if err == nil {
		return out, nil
	}
	if apiErr, ok := err.(*APIResponseError); !ok || apiErr.StatusCode != http.StatusNotFound {
		return out, err
	}

	return out, c.doJSON(ctx, http.MethodPost, "/templates", nil, body, &out)
}

func (c *Client) CreateTemplateV2(ctx context.Context, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	// Cube compat: /v2/templates is pure-404 on Cube, use legacy route directly.
	if c.compatMode {
		return out, c.doJSON(ctx, http.MethodPost, "/templates", nil, body, &out)
	}

	// E2B native: try /v2 → /templates fallback chain.
	err := c.doJSON(ctx, http.MethodPost, "/v2/templates", nil, body, &out)
	if err == nil {
		return out, nil
	}
	if apiErr, ok := err.(*APIResponseError); !ok || apiErr.StatusCode != http.StatusNotFound {
		return out, err
	}

	return out, c.doJSON(ctx, http.MethodPost, "/templates", nil, body, &out)
}

func (c *Client) GetBuildUploadLink(ctx context.Context, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, "/templates/upload-link", query, nil, &out)
	return out, err
}

func (c *Client) ListTemplates(ctx context.Context, query url.Values) ([]Template, error) {
	var out []Template
	err := c.doJSON(ctx, http.MethodGet, "/templates", query, nil, &out)
	return out, err
}

func (c *Client) CreateTemplate(ctx context.Context, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPost, "/templates", nil, body, &out)
	return out, err
}

func (c *Client) GetTemplate(ctx context.Context, templateID string, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/templates/%s", templateID), query, nil, &out)
	return out, err
}

func (c *Client) RebuildTemplate(ctx context.Context, templateID string, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/templates/%s", templateID), nil, body, &out)
	return out, err
}

func (c *Client) DeleteTemplate(ctx context.Context, templateID string) error {
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/templates/%s", templateID), nil, nil, nil)
}

func (c *Client) UpdateTemplate(ctx context.Context, templateID string, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPatch, fmt.Sprintf("/templates/%s", templateID), nil, body, &out)
	return out, err
}

func (c *Client) StartBuild(ctx context.Context, templateID string, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/templates/%s/build", templateID), nil, body, &out)
	return out, err
}

func (c *Client) StartBuildV2(ctx context.Context, templateID string, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/v2/templates/%s/build", templateID), nil, body, &out)
	return out, err
}

func (c *Client) UpdateTemplateV2(ctx context.Context, templateID string, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPatch, fmt.Sprintf("/v2/templates/%s", templateID), nil, body, &out)
	return out, err
}

func (c *Client) GetBuildStatus(ctx context.Context, templateID, buildID string) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/templates/%s/builds/%s", templateID, buildID), nil, nil, &out)
	return out, err
}

func (c *Client) GetBuildLogs(ctx context.Context, templateID, buildID string, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/templates/%s/builds/%s/logs", templateID, buildID), query, nil, &out)
	return out, err
}

func (c *Client) GetTemplateByAlias(ctx context.Context, alias string) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/templates/aliases/%s", alias), nil, nil, &out)
	return out, err
}
