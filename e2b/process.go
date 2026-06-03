package e2b

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
	Cmd         string            `json:"cmd,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Args        []string          `json:"args,omitempty"`
	PTY         bool              `json:"pty,omitempty"`
	SandboxID   string            `json:"sandboxID,omitempty"`
	ContainerID string            `json:"containerID,omitempty"`
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
	// Cube compat: /process/start is pure-404 on Cube, go directly to compat route.
	if c.compatMode {
		return c.startProcessCompat(ctx, req)
	}

	// E2B native: try /process/start, fall back to /sandbox/exec on 404.
	var nativeOut StartProcessResponse
	err := c.doJSON(ctx, http.MethodPost, "/process/start", nil, req, &nativeOut)
	if err == nil {
		if nativeOut.PID == "" {
			nativeOut.PID = "process-start"
		}
		return &nativeOut, nil
	}
	if apiErr, ok := err.(*APIResponseError); !ok || apiErr.StatusCode != http.StatusNotFound {
		return nil, err
	}

	return c.startProcessCompat(ctx, req)
}

// startProcessCompat uses the Cube-compatible /sandbox/exec route.
func (c *Client) startProcessCompat(ctx context.Context, req StartProcessRequest) (*StartProcessResponse, error) {
	sandboxID := req.SandboxID
	if sandboxID == "" && req.Env != nil {
		sandboxID = req.Env["E2B_SANDBOX_ID"]
	}
	if sandboxID == "" {
		return nil, &BaseError{Message: "missing sandbox id: set StartProcessRequest.SandboxID or env E2B_SANDBOX_ID"}
	}

	containerID := req.ContainerID
	if containerID == "" {
		containerID = sandboxID
	}

	args := req.Args
	if len(args) == 0 && req.Cmd != "" {
		args = strings.Fields(req.Cmd)
	}
	if len(args) == 0 {
		return nil, &BaseError{Message: "missing process args: set StartProcessRequest.Args or Cmd"}
	}

	payload := map[string]interface{}{
		"sandbox_id":   sandboxID,
		"container_id": containerID,
		"args":         args,
	}
	var out StartProcessResponse
	err := c.doJSON(ctx, http.MethodPost, "/sandbox/exec", nil, payload, &out)
	if err == nil {
		if out.PID == "" {
			out.PID = "cube-exec"
		}
		return &out, nil
	}
	return nil, err
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
