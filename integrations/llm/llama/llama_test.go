package llama

import (
	"context"
	"os"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("", "")
	if c.Model != "llama3.1:8b" {
		t.Errorf("default model = %q, want llama3.1:8b", c.Model)
	}
	if c.BaseURL != "http://localhost:11434" {
		t.Errorf("default base URL = %q, want http://localhost:11434", c.BaseURL)
	}
}

func TestNewClient_Custom(t *testing.T) {
	c := NewClient("llama3.2:3b", "http://remote:11434")
	if c.Model != "llama3.2:3b" {
		t.Errorf("model = %q, want llama3.2:3b", c.Model)
	}
	if c.BaseURL != "http://remote:11434" {
		t.Errorf("base URL = %q, want http://remote:11434", c.BaseURL)
	}
}

func TestComplete_Integration(t *testing.T) {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		t.Skip("OLLAMA_HOST not set, skipping integration test")
	}

	c := NewClient("llama3.1:8b", host)
	resp, err := c.Complete(context.Background(), "Say hello in one word.", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp == "" {
		t.Error("expected non-empty response")
	}
}
