package e2b

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// const defaultBaseURL = "https://api.e2b.app"
const defaultBaseURL = "http://127.0.0.1:13000"

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	compatMode bool
}

type ClientOption func(*Client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

func WithCompatMode(enabled bool) ClientOption {
	return func(c *Client) {
		c.compatMode = enabled
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		apiKey:  "API_KEY",
		// Default on for local CubeSandbox compatibility.
		compatMode: true,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type APIResponseError struct {
	StatusCode int
	Body       string
}

func (e *APIResponseError) Error() string {
	return fmt.Sprintf("e2b api error: status=%d body=%s", e.StatusCode, e.Body)
}

func (c *Client) doJSON(ctx context.Context, method, path string, query url.Values, requestBody any, out any) error {
	if c == nil {
		return &BaseError{Message: "nil client"}
	}

	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var body io.Reader
	if requestBody != nil {
		b, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		fmt.Println("doJSON:", string(b))
		body = bytes.NewReader(b)
	}

	fmt.Println("doJSON:", method, u)
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return err
	}

	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var payload APIErrorPayload
		if json.Unmarshal(respBytes, &payload) == nil {
			if payload.Message != "" || payload.Type != "" || payload.Code != "" {
				return &BaseError{Type: ErrorType(payload.Type), Message: payload.Message}
			}
		}
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
