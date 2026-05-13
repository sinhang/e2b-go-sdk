package e2b

import "testing"

func TestNewClientDefaults(t *testing.T) {
	c := NewClient(WithAPIKey("api_key"))
	if c.baseURL != defaultBaseURL {
		t.Fatalf("unexpected baseURL: %s", c.baseURL)
	}
	if c.apiKey != "api_key" {
		t.Fatalf("unexpected api key")
	}
}
