package e2b

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// codeInterpreterPort is the port the lightweight code interpreter listens on.
const codeInterpreterPort = 49999

// ---- Request / Response types ----

// RunCodeRequest is the request body for executing code.
type RunCodeRequest struct {
	SandboxID string `json:"sandboxID,omitempty"`
	Code      string `json:"code"`
	Language  string `json:"language,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
	ContextID string `json:"contextID,omitempty"`
}

// Execution holds the aggregated result of a code execution.
type Execution struct {
	Results []Result        `json:"results,omitempty"`
	Logs    Logs            `json:"logs,omitempty"`
	Error   *ExecutionError `json:"error,omitempty"`
}

// Text returns a concatenated string of all text results.
func (e *Execution) Text() string {
	var b strings.Builder
	for _, r := range e.Results {
		if r.Text != "" {
			b.WriteString(r.Text)
		}
	}
	return b.String()
}

// Result represents a single output chunk (stdout, stderr, etc.).
type Result struct {
	Text      string `json:"text,omitempty"`
	Type      string `json:"type,omitempty"`
	IsStderr  bool   `json:"isStderr,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// Logs holds stdout/stderr collected during execution.
type Logs struct {
	Stdout []string `json:"stdout,omitempty"`
	Stderr []string `json:"stderr,omitempty"`
}

// OutputMessage is a single NDJSON line from the /execute endpoint.
type OutputMessage struct {
	Type           string `json:"type,omitempty"`
	Text           string `json:"text,omitempty"`
	Timestamp      string `json:"timestamp,omitempty"`
	ExecutionCount int    `json:"execution_count,omitempty"`
	IsStderr       bool   `json:"is_stderr,omitempty"`
}

// ExecutionError holds error information from code execution.
type ExecutionError struct {
	Name      string `json:"name,omitempty"`
	Value     string `json:"value,omitempty"`
	Traceback string `json:"traceback,omitempty"`
}

// CodeContext represents a persistent execution context.
type CodeContext struct {
	ID       string `json:"id,omitempty"`
	Language string `json:"language,omitempty"`
}

// CreateCodeContextRequest is the request to create a code context.
type CreateCodeContextRequest struct {
	SandboxID string `json:"sandboxID,omitempty"`
	Language  string `json:"language,omitempty"`
}

// RemoveCodeContextRequest is the request to remove a code context.
type RemoveCodeContextRequest struct {
	SandboxID string `json:"sandboxID,omitempty"`
	ContextID string `json:"contextID,omitempty"`
}

// RestartCodeContextRequest is the request to restart a code context.
type RestartCodeContextRequest struct {
	SandboxID string `json:"sandboxID,omitempty"`
	ContextID string `json:"contextID,omitempty"`
}

// ---- CodeInterpreter ----

// CodeInterpreter executes code inside a sandbox through the data plane.
// It connects to the lightweight Python code interpreter service
// running on port 49999 inside the sandbox.
type CodeInterpreter struct {
	router *SandboxRouter
}

// Run executes code and returns the aggregated execution result.
func (ci *CodeInterpreter) Run(ctx context.Context, req RunCodeRequest) (*Execution, error) {
	body := map[string]interface{}{
		"code": req.Code,
	}
	if req.Language != "" {
		body["language"] = req.Language
	} else {
		body["language"] = "python"
	}
	if req.ContextID != "" {
		body["context_id"] = req.ContextID
	}

	respBytes, err := ci.router.doRaw(ctx, "POST", codeInterpreterPort, "/execute", nil, body)
	if err != nil {
		return nil, err
	}

	return parseExecuteResponse(respBytes)
}

// RunSimple is a convenience that executes a code string as Python.
func (ci *CodeInterpreter) RunSimple(ctx context.Context, code string) (*Execution, error) {
	return ci.Run(ctx, RunCodeRequest{Code: code, Language: "python"})
}

// parseExecuteResponse parses the NDJSON (newline-delimited JSON) response
// from the lightweight code interpreter's /execute endpoint.
//
// Example response:
//
//	{"type": "number_of_executions", "execution_count": 1}
//	{"type": "stdout", "text": "2\n", "timestamp": "..."}
//	{"type": "end_of_execution"}
func parseExecuteResponse(data []byte) (*Execution, error) {
	exec := &Execution{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var msg OutputMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "stdout":
			exec.Results = append(exec.Results, Result{
				Text:      msg.Text,
				Type:      "stdout",
				Timestamp: msg.Timestamp,
			})
			exec.Logs.Stdout = append(exec.Logs.Stdout, msg.Text)
		case "stderr":
			exec.Results = append(exec.Results, Result{
				Text:      msg.Text,
				Type:      "stderr",
				IsStderr:  true,
				Timestamp: msg.Timestamp,
			})
			exec.Logs.Stderr = append(exec.Logs.Stderr, msg.Text)
		case "error":
			exec.Error = &ExecutionError{
				Name:      msg.Type,
				Value:     msg.Text,
				Traceback: msg.Text,
			}
		case "number_of_executions", "end_of_execution":
			// metadata, no action needed
		default:
			// capture unknown types as text results
			if msg.Text != "" {
				exec.Results = append(exec.Results, Result{
					Text:      msg.Text,
					Type:      msg.Type,
					Timestamp: msg.Timestamp,
				})
			}
		}
	}
	return exec, scanner.Err()
}

// ---- Context management (stubs for future use) ----

// CreateCodeContext creates a persistent code execution context.
func (ci *CodeInterpreter) CreateCodeContext(ctx context.Context, req CreateCodeContextRequest) (*CodeContext, error) {
	payload := map[string]interface{}{}
	if req.Language != "" {
		payload["language"] = req.Language
	}
	var out CodeContext
	err := ci.router.doJSON(ctx, "POST", codeInterpreterPort, "/contexts", nil, payload, &out)
	return &out, err
}

// ListCodeContexts returns all code contexts in a sandbox.
func (ci *CodeInterpreter) ListCodeContexts(ctx context.Context) ([]CodeContext, error) {
	var out []CodeContext
	err := ci.router.doJSON(ctx, "GET", codeInterpreterPort, "/contexts", nil, nil, &out)
	return out, err
}

// RemoveCodeContext deletes a code context.
func (ci *CodeInterpreter) RemoveCodeContext(ctx context.Context, req RemoveCodeContextRequest) error {
	path := fmt.Sprintf("/contexts/%s", req.ContextID)
	return ci.router.doJSON(ctx, "DELETE", codeInterpreterPort, path, nil, nil, nil)
}

// RestartCodeContext restarts a code context.
func (ci *CodeInterpreter) RestartCodeContext(ctx context.Context, req RestartCodeContextRequest) error {
	path := fmt.Sprintf("/contexts/%s/restart", req.ContextID)
	return ci.router.doJSON(ctx, "POST", codeInterpreterPort, path, nil, nil, nil)
}

// ---- Client accessors ----

// NewCodeInterpreter creates a CodeInterpreter for the given sandbox,
// using the client's data-plane configuration for routing.
func NewCodeInterpreter(client *Client, sandboxID string) *CodeInterpreter {
	return &CodeInterpreter{
		router: &SandboxRouter{
			sandboxID:    sandboxID,
			domain:       client.sandboxDomain,
			dataPlaneURL: client.dataPlaneURL,
			httpClient:   client.dataPlaneClient,
		},
	}
}

// ExecuteCode is a standalone convenience that executes code in one call.
func ExecuteCode(ctx context.Context, client *Client, code string, sandboxID string) (*Execution, error) {
	return NewCodeInterpreter(client, sandboxID).Run(ctx, RunCodeRequest{
		SandboxID: sandboxID,
		Code:      code,
		Language:  "python",
	})
}
