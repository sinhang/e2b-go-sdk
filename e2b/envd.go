package e2b

import (
	"context"
	"net/http"
)

type EnvdHealthResponse struct {
	Status string `json:"status,omitempty"`
}

type EnvdStatsResponse struct {
	CPUCount int64                  `json:"cpuCount,omitempty"`
	MemoryMB int64                  `json:"memoryMB,omitempty"`
	Raw      map[string]interface{} `json:"-"`
}

type EnvdEnvironmentVariablesResponse struct {
	Variables map[string]string `json:"variables,omitempty"`
}

func (c *Client) EnvdHealth(ctx context.Context) (*EnvdHealthResponse, error) {
	var out EnvdHealthResponse
	err := c.doJSON(ctx, http.MethodGet, "/envd/health", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) EnvdStats(ctx context.Context) (*EnvdStatsResponse, error) {
	var out EnvdStatsResponse
	err := c.doJSON(ctx, http.MethodGet, "/envd/stats", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) EnvdEnvs(ctx context.Context) (*EnvdEnvironmentVariablesResponse, error) {
	var out EnvdEnvironmentVariablesResponse
	err := c.doJSON(ctx, http.MethodGet, "/envd/envs", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
