package e2b

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CreateSandboxRequest struct {
	TemplateID          string    `json:"templateID,omitempty"`
	Timeout             int       `json:"timeout,omitempty"`
	AutoPause           *bool     `json:"autoPause,omitempty"`
	Secure              *bool     `json:"secure,omitempty"`
	AllowInternetAccess *bool     `json:"allow_internet_access,omitempty"`
	Network             JSONMap   `json:"network,omitempty"`
	Metadata            JSONMap   `json:"metadata,omitempty"`
	EnvVars             JSONMap   `json:"envVars,omitempty"`
	MCP                 JSONMap   `json:"mcp,omitempty"`
	VolumeMounts        []JSONMap `json:"volumeMounts,omitempty"`
}

type SnapshotRequest struct {
	Name string `json:"name,omitempty"`
}

func (c *Client) ListSandboxes(ctx context.Context, metadata string) ([]Sandbox, error) {
	q := url.Values{}
	if metadata != "" {
		q.Set("metadata", metadata)
	}
	var out []Sandbox
	err := c.doJSON(ctx, http.MethodGet, "/sandboxes", q, nil, &out)
	return out, err
}

func (c *Client) CreateSandbox(ctx context.Context, req CreateSandboxRequest) (*Sandbox, error) {
	var out Sandbox
	err := c.doJSON(ctx, http.MethodPost, "/sandboxes", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListSandboxesV2(ctx context.Context, metadata string) ([]Sandbox, error) {
	q := url.Values{}
	if metadata != "" {
		q.Set("metadata", metadata)
	}
	var out []Sandbox
	err := c.doJSON(ctx, http.MethodGet, "/v2/sandboxes", q, nil, &out)
	return out, err
}

func (c *Client) ListSandboxMetrics(ctx context.Context, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, "/sandboxes/metrics", query, nil, &out)
	return out, err
}

func (c *Client) GetSandboxLogs(ctx context.Context, sandboxID string, start int64, limit int) (JSONMap, error) {
	q := url.Values{}
	if start > 0 {
		q.Set("start", strconv.FormatInt(start, 10))
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/sandboxes/%s/logs", sandboxID), q, nil, &out)
	return out, err
}

func (c *Client) GetSandboxLogsV2(ctx context.Context, sandboxID string, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/v2/sandboxes/%s/logs", sandboxID), query, nil, &out)
	return out, err
}

func (c *Client) GetSandbox(ctx context.Context, sandboxID string) (*Sandbox, error) {
	var out Sandbox
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/sandboxes/%s", sandboxID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSandbox(ctx context.Context, sandboxID string) error {
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/sandboxes/%s", sandboxID), nil, nil, nil)
}

func (c *Client) GetSandboxMetrics(ctx context.Context, sandboxID string, query url.Values) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/sandboxes/%s/metrics", sandboxID), query, nil, &out)
	return out, err
}

func (c *Client) PauseSandbox(ctx context.Context, sandboxID string) (*Sandbox, error) {
	var out Sandbox
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/sandboxes/%s/pause", sandboxID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ResumeSandbox(ctx context.Context, sandboxID string) (*Sandbox, error) {
	var out Sandbox
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/sandboxes/%s/resume", sandboxID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ConnectToSandbox(ctx context.Context, sandboxID string, body JSONMap) (*Sandbox, error) {
	var out Sandbox
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/sandboxes/%s/connect", sandboxID), nil, body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SetSandboxTimeout(ctx context.Context, sandboxID string, timeoutSec int) (*Sandbox, error) {
	payload := JSONMap{"timeout": timeoutSec}
	var out Sandbox
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/sandboxes/%s/timeout", sandboxID), nil, payload, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PutSandboxesNetwork(ctx context.Context, sandboxID string, body JSONMap) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPut, fmt.Sprintf("/sandboxes/%s/network", sandboxID), nil, body, &out)
	return out, err
}

func (c *Client) RefreshSandbox(ctx context.Context, sandboxID string, timeoutSec int) (*Sandbox, error) {
	payload := JSONMap{}
	if timeoutSec > 0 {
		payload["timeout"] = timeoutSec
	}
	var out Sandbox
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/sandboxes/%s/refresh", sandboxID), nil, payload, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PostSandboxesSnapshots(ctx context.Context, sandboxID string, req SnapshotRequest) (JSONMap, error) {
	out := JSONMap{}
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/sandboxes/%s/snapshots", sandboxID), nil, req, &out)
	return out, err
}

func (c *Client) GetSnapshots(ctx context.Context, query url.Values) ([]JSONMap, error) {
	var out []JSONMap
	err := c.doJSON(ctx, http.MethodGet, "/snapshots", query, nil, &out)
	return out, err
}
