package e2b

import (
	"context"
	"net/http"
	"strings"
)

// RunCommandRequest is the full request for running a command inside a sandbox.
type RunCommandRequest struct {
	SandboxID string            `json:"sandboxID,omitempty"`
	Cmd       string            `json:"cmd,omitempty"`
	Args      []string          `json:"args,omitempty"`
	EnvVars   map[string]string `json:"envVars,omitempty"`
	Timeout   int               `json:"timeout,omitempty"`
}

// CommandResult holds the output of a command execution.
type CommandResult struct {
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	ExitCode int    `json:"exitCode,omitempty"`
	PID      string `json:"pid,omitempty"`
}

// StdoutText returns the stdout content trimmed of trailing newline.
func (r *CommandResult) StdoutText() string {
	result := r.Stdout
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// Commands provides methods for executing shell commands inside a sandbox.
// Execution goes through the control plane (CubeAPI) to envd inside the sandbox.
type Commands struct {
	client    *Client
	sandboxID string
}

// Run executes a command inside the sandbox and returns its output.
//
// The command is sent to CubeAPI via the control plane, which proxies to
// envd inside the target sandbox. On success the response includes stdout,
// stderr, exit code, and PID.
func (cmd *Commands) Run(ctx context.Context, req RunCommandRequest) (*CommandResult, error) {
	if req.SandboxID == "" {
		req.SandboxID = cmd.sandboxID
	}
	if req.SandboxID == "" {
		return nil, &BaseError{Message: "missing sandbox id"}
	}

	args := req.Args
	if len(args) == 0 && req.Cmd != "" {
		args = strings.Fields(req.Cmd)
	}
	if len(args) == 0 {
		return nil, &BaseError{Message: "missing command: set Cmd or Args"}
	}

	// 1. Try E2B-native /process/start (verified working on Cube).
	e2bReq := StartProcessRequest{
		Cmd:       req.Cmd,
		Args:      req.Args,
		Env:       req.EnvVars,
		SandboxID: req.SandboxID,
	}
	var e2bOut StartProcessResponse
	err := cmd.client.doJSON(ctx, http.MethodPost, "/process/start", nil, e2bReq, &e2bOut)
	if err == nil {
		return &CommandResult{PID: e2bOut.PID}, nil
	}

	// 2. Try CubeSandbox compat format with the same path.
	snakePayload := map[string]interface{}{
		"sandbox_id":   req.SandboxID,
		"container_id": req.SandboxID,
		"args":         args,
	}
	var raw map[string]interface{}
	err2 := cmd.client.doJSON(ctx, http.MethodPost, "/process/start", nil, snakePayload, &raw)
	if err2 == nil {
		return parseCommandResult(raw), nil
	}

	return nil, err
}

// parseCommandResult extracts stdout/stderr/exitCode from a loose response shape.
func parseCommandResult(raw map[string]interface{}) *CommandResult {
	r := &CommandResult{PID: "cube-exec"}
	if v, ok := raw["stdout"].(string); ok {
		r.Stdout = v
	}
	if v, ok := raw["stderr"].(string); ok {
		r.Stderr = v
	}
	if v, ok := raw["exit_code"].(float64); ok {
		r.ExitCode = int(v)
	} else if v, ok := raw["exitCode"].(float64); ok {
		r.ExitCode = int(v)
	}
	if ret, ok := raw["ret"].(map[string]interface{}); ok {
		if msg, ok := ret["ret_msg"].(string); ok && r.Stdout == "" {
			r.Stdout = msg
		}
		if code, ok := ret["ret_code"].(float64); ok && r.ExitCode == 0 {
			r.ExitCode = int(code)
		}
	}
	return r
}

// RunSimple is a convenience that runs a command string on a sandbox.
func (cmd *Commands) RunSimple(ctx context.Context, command string) (*CommandResult, error) {
	return cmd.Run(ctx, RunCommandRequest{
		SandboxID: cmd.sandboxID,
		Cmd:       command,
	})
}

// Commands returns a Commands runner bound to the given sandbox.
// Command execution uses the control plane (CubeAPI).
func (c *Client) Commands(sandboxID string) *Commands {
	return &Commands{
		client:    c,
		sandboxID: sandboxID,
	}
}

// RunCommand is a standalone convenience that creates a Commands runner
// and executes the request in one call.
func (c *Client) RunCommand(ctx context.Context, req RunCommandRequest) (*CommandResult, error) {
	return c.Commands(req.SandboxID).Run(ctx, req)
}
