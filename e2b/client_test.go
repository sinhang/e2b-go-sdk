package e2b

import "testing"

func TestNewClientDefaults(t *testing.T) {
	c := NewClient("k")
	if c.baseURL != defaultBaseURL {
		t.Fatalf("unexpected baseURL: %s", c.baseURL)
	}
	if c.apiKey != "k" {
		t.Fatalf("unexpected api key")
	}
}
