package e2b

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// SandboxRouter constructs data-plane URLs for accessing services
// (envd, Jupyter, files) inside a sandbox.
//
// Dev mode (dataPlaneURL set):
//
//	URL:  {dataPlaneURL}/path
//	Host: {port}-{sandboxID}.{domain}
//
// Production mode (dataPlaneURL empty):
//
//	URL:  https://{port}-{sandboxID}.{domain}/path
//	Host: (derived from URL by http.Client)
type SandboxRouter struct {
	sandboxID    string
	domain       string
	dataPlaneURL string
	httpClient   *http.Client
}

// BuildHost returns the Host header value CubeProxy uses for routing.
//
//	{port}-{sandboxID}.{domain}
//
// Example: "49983-abc123def456.cube.app"
func (r *SandboxRouter) BuildHost(port int) string {
	return fmt.Sprintf("%d-%s.%s", port, r.sandboxID, r.domain)
}

// BuildURL constructs the full URL for accessing a sandbox service.
// In dev mode the target is dataPlaneURL + path with a synthetic Host header.
// In production the target URL embeds the hostname so DNS resolves it.
func (r *SandboxRouter) BuildURL(port int, path string) string {
	if r.dataPlaneURL != "" {
		return fmt.Sprintf("%s%s", r.dataPlaneURL, path)
	}
	return fmt.Sprintf("https://%s%s", r.BuildHost(port), path)
}

// doJSON sends a JSON request through the data plane and unmarshals the response.
func (r *SandboxRouter) doJSON(ctx context.Context, method string, port int, path string, query url.Values, requestBody any, out any) error {
	u := r.BuildURL(port, path)
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var body io.Reader
	if requestBody != nil {
		b, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		fmt.Println("dp doJSON:", string(b))
		body = bytes.NewReader(b)
	}

	fmt.Println("dp doJSON:", method, u)
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return err
	}

	// Always set the Host header for CubeProxy routing.
	req.Host = r.BuildHost(port)

	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIResponseError{StatusCode: resp.StatusCode, Body: string(respBytes)}
	}

	if out == nil || len(respBytes) == 0 {
		return nil
	}

	if err := json.Unmarshal(respBytes, out); err != nil {
		return err
	}
	return nil
}

// doRaw sends a raw request through the data plane and returns the response body.
func (r *SandboxRouter) doRaw(ctx context.Context, method string, port int, path string, query url.Values, requestBody any) ([]byte, error) {
	u := r.BuildURL(port, path)
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var body io.Reader
	if requestBody != nil {
		b, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		fmt.Println("dp doRaw:", string(b))
		body = bytes.NewReader(b)
	}

	fmt.Println("dp doRaw:", method, u)
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}

	req.Host = r.BuildHost(port)

	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIResponseError{StatusCode: resp.StatusCode, Body: string(respBytes)}
	}

	return respBytes, nil
}
