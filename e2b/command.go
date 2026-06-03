package e2b

import (
	"context"
	"net/http"
	"strings"
)

// envdPort is the port envd listens on inside the sandbox for process/command operations.
const envdPort = 49983

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
type Commands struct {
	client    *Client
	sandboxID string
	router    *SandboxRouter // lazy init for data plane access

	// EnvdPort is the port envd listens on inside the sandbox.
	// Default is 49983; set before calling Run if your template uses a different port.
	EnvdPort int
}

func (cmd *Commands) getRouter() *SandboxRouter {
	if cmd.router == nil {
		cmd.router = &SandboxRouter{
			sandboxID:    cmd.sandboxID,
			domain:       cmd.client.sandboxDomain,
			dataPlaneURL: cmd.client.dataPlaneURL,
			httpClient:   cmd.client.dataPlaneClient,
		}
	}
	return cmd.router
}

// Run executes a command inside the sandbox and returns its output.
//
// In compat mode (Cube): executes through the data plane (CubeProxy → envd
// inside the sandbox). In non-compat mode (E2B): uses the control plane.
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

	// Cube compat: use data plane to reach envd inside the sandbox.
	if cmd.client.compatMode {
		return cmd.runViaDataPlane(ctx, req, args)
	}

	// E2B native: use control plane /process/start.
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

	return nil, err
}

// runViaDataPlane sends the command to envd inside the sandbox through the
// data plane (CubeProxy → sandbox envd port).
func (cmd *Commands) runViaDataPlane(ctx context.Context, req RunCommandRequest, args []string) (*CommandResult, error) {
	router := cmd.getRouter()

	port := cmd.EnvdPort
	if port == 0 {
		port = envdPort
	}

	payload := map[string]interface{}{
		"cmd":  req.Cmd,
		"args": args,
	}
	if req.EnvVars != nil {
		payload["env"] = req.EnvVars
	}

	var raw map[string]interface{}
	err := router.doJSON(ctx, http.MethodPost, port, "/process/start", nil, payload, &raw)
	if err != nil {
		return nil, err
	}
	return parseCommandResult(raw), nil
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
