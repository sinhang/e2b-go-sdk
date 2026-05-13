package e2b

import (
	"context"
	"fmt"
	"net/http"
)

type CloseStdinRequest struct {
	PID string `json:"pid,omitempty"`
}
type CloseStdinResponse struct{}

type ConnectProcessRequest struct {
	PID string `json:"pid,omitempty"`
}
type ConnectProcessResponse struct {
	PID string `json:"pid,omitempty"`
}

type ListProcessRequest struct {
	IncludeExited bool `json:"includeExited,omitempty"`
}
type ProcessInfo struct {
	PID    string `json:"pid,omitempty"`
	Status string `json:"status,omitempty"`
	Cmd    string `json:"cmd,omitempty"`
}
type ListProcessResponse struct {
	Processes []ProcessInfo `json:"processes,omitempty"`
}

type SendInputRequest struct {
	PID   string `json:"pid,omitempty"`
	Input string `json:"input,omitempty"`
}
type SendInputResponse struct{}

type SendSignalRequest struct {
	PID    string `json:"pid,omitempty"`
	Signal string `json:"signal,omitempty"`
}
type SendSignalResponse struct{}

type StartProcessRequest struct {
	Cmd  string            `json:"cmd,omitempty"`
	Env  map[string]string `json:"env,omitempty"`
	Args []string          `json:"args,omitempty"`
	PTY  bool              `json:"pty,omitempty"`
}
type StartProcessResponse struct {
	PID string `json:"pid,omitempty"`
}

type StreamInputRequest struct {
	PID   string `json:"pid,omitempty"`
	Input string `json:"input,omitempty"`
}
type StreamInputResponse struct{}

type UpdateProcessRequest struct {
	PID string `json:"pid,omitempty"`
}
type UpdateProcessResponse struct {
	PID    string `json:"pid,omitempty"`
	Status string `json:"status,omitempty"`
}

func (c *Client) CloseStdin(ctx context.Context, req CloseStdinRequest) (*CloseStdinResponse, error) {
	var out CloseStdinResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/closestdin", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ConnectProcess(ctx context.Context, req ConnectProcessRequest) (*ConnectProcessResponse, error) {
	var out ConnectProcessResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/connect", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListProcess(ctx context.Context, req ListProcessRequest) (*ListProcessResponse, error) {
	var out ListProcessResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/list", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendInput(ctx context.Context, req SendInputRequest) (*SendInputResponse, error) {
	var out SendInputResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/sendinput", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendSignal(ctx context.Context, req SendSignalRequest) (*SendSignalResponse, error) {
	var out SendSignalResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/sendsignal", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) StartProcess(ctx context.Context, req StartProcessRequest) (*StartProcessResponse, error) {
	var out StartProcessResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/start", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) StreamInput(ctx context.Context, req StreamInputRequest) (*StreamInputResponse, error) {
	var out StreamInputResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/streaminput", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateProcess(ctx context.Context, req UpdateProcessRequest) (*UpdateProcessResponse, error) {
	var out UpdateProcessResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/update", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ProcessRPC(ctx context.Context, method string, body interface{}) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/process/%s", method), nil, body, &out)
	return out, err
}
