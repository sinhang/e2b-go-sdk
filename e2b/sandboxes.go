package e2b

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type cubeCreateSandboxResponse struct {
	RequestID string `json:"RequestID,omitempty"`
	SandboxID string `json:"sandbox_id,omitempty"`
	SandboxIP string `json:"sandbox_ip,omitempty"`
	HostID    string `json:"host_id,omitempty"`
	HostIP    string `json:"host_ip,omitempty"`
	State     string `json:"state,omitempty"`
	Ret       struct {
		RetCode int    `json:"ret_code,omitempty"`
		RetMsg  string `json:"ret_msg,omitempty"`
	} `json:"ret,omitempty"`
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
	if err == nil {
		return &out, nil
	}

	// CubeSandbox compatibility:
	// Some deployments don't expose POST /sandboxes and only provide /v2/sandboxes
	// with snake_case fields. Retry on 404 with alternate shapes.
	var apiErr *APIResponseError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound {
		return nil, err
	}

	altReq := map[string]interface{}{
		"template_id": req.TemplateID,
	}
	if req.Timeout > 0 {
		altReq["timeout"] = req.Timeout
	}
	if req.AutoPause != nil {
		altReq["auto_pause"] = *req.AutoPause
	}
	if req.Secure != nil {
		altReq["secure"] = *req.Secure
	}
	if req.AllowInternetAccess != nil {
		altReq["allow_internet_access"] = *req.AllowInternetAccess
	}
	if req.Network != nil {
		altReq["network"] = req.Network
	}
	if req.Metadata != nil {
		altReq["metadata"] = req.Metadata
	}
	if req.EnvVars != nil {
		altReq["env_vars"] = req.EnvVars
	}
	if req.MCP != nil {
		altReq["mcp"] = req.MCP
	}
	if len(req.VolumeMounts) > 0 {
		altReq["volume_mounts"] = req.VolumeMounts
	}

	err = c.doJSON(ctx, http.MethodPost, "/v2/sandboxes", nil, altReq, &out)
	if err == nil {
		return &out, nil
	}

	// CubeMaster compatibility path.
	var cubeResp cubeCreateSandboxResponse
	err = c.doJSON(ctx, http.MethodPost, "/cube/sandbox", nil, altReq, &cubeResp)
	if err == nil && cubeResp.SandboxID != "" {
		return &Sandbox{
			SandboxID: cubeResp.SandboxID,
			State:     cubeResp.State,
		}, nil
	}

	// Some installations expose CubeMaster on :8089 while CubeAPI is on :3002.
	altBase := convertPort(c.baseURL, "8089")
	if altBase != "" && altBase != c.baseURL {
		altClient := *c
		altClient.baseURL = altBase
		err = altClient.doJSON(ctx, http.MethodPost, "/cube/sandbox", nil, altReq, &cubeResp)
		if err == nil && cubeResp.SandboxID != "" {
			return &Sandbox{
				SandboxID: cubeResp.SandboxID,
				State:     cubeResp.State,
			}, nil
		}
	}
	return nil, err
}

func convertPort(raw, port string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return ""
	}
	host := u.Hostname()
	if host == "" {
		return ""
	}
	u.Host = host + ":" + port
	return strings.TrimRight(u.String(), "/")
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
