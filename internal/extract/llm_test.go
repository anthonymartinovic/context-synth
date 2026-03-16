package extract

import (
	"context"
	"fmt"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type mockClient struct {
	response string
	err      error
}

func (m *mockClient) Complete(_ context.Context, _, _ string) (string, error) {
	return m.response, m.err
}

func TestLLMExtractor_ClassifiesIntoSections(t *testing.T) {
	client := &mockClient{
		response: `[
			{"section": "Domain Knowledge", "content": "The platform handles products and orders."},
			{"section": "Architecture Decisions", "content": "PostgreSQL over DynamoDB for consistency."}
		]`,
	}

	ext := &LLMExtractor{Client: client}
	snap := protocol.Snapshot{
		Items: []protocol.SnapshotItem{
			{
				Source: protocol.Source{
					Path:        "docs/overview.md",
					Content:     []byte("test content"),
					ContentHash: "abc123def456ghij",
					TokenCount:  100,
					Weight:      1.0,
					DeclOrder:   0,
				},
			},
		},
	}

	sections := []config.SectionDecl{
		{Name: "Domain Knowledge", Budget: 0.5},
		{Name: "Architecture Decisions", Budget: 0.5},
	}

	results, err := ext.Extract(context.Background(), snap, sections)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d extractions, want 2", len(results))
	}
	if results[0].Section != "Domain Knowledge" {
		t.Errorf("first section = %q, want Domain Knowledge", results[0].Section)
	}
	if results[1].Section != "Architecture Decisions" {
		t.Errorf("second section = %q, want Architecture Decisions", results[1].Section)
	}
}

func TestLLMExtractor_InheritsProvenance(t *testing.T) {
	client := &mockClient{
		response: `[{"section": "Overview", "content": "Some knowledge item."}]`,
	}

	ext := &LLMExtractor{Client: client}
	snap := protocol.Snapshot{
		Items: []protocol.SnapshotItem{
			{
				Source: protocol.Source{
					Path:        "docs/test.md",
					Content:     []byte("source text"),
					ContentHash: "hashvalue1234abcd",
					TokenCount:  50,
					Weight:      0.8,
					DeclOrder:   3,
				},
			},
		},
	}

	sections := []config.SectionDecl{{Name: "Overview", Budget: 1.0}}
	results, err := ext.Extract(context.Background(), snap, sections)
	if err != nil {
		t.Fatal(err)
	}

	if results[0].Weight != 0.8 {
		t.Errorf("weight = %f, want 0.8", results[0].Weight)
	}
	if results[0].SourcePath != "docs/test.md" {
		t.Errorf("source = %q, want docs/test.md", results[0].SourcePath)
	}
	if results[0].ContentHash != "hashvalue1234abcd" {
		t.Errorf("hash = %q, want hashvalue1234abcd", results[0].ContentHash)
	}
	if results[0].DeclOrder != 3 {
		t.Errorf("decl order = %d, want 3", results[0].DeclOrder)
	}
}

func TestLLMExtractor_HandlesCodeFences(t *testing.T) {
	client := &mockClient{
		response: "```json\n[{\"section\": \"Info\", \"content\": \"fenced response\"}]\n```",
	}

	ext := &LLMExtractor{Client: client}
	snap := protocol.Snapshot{
		Items: []protocol.SnapshotItem{
			{
				Source: protocol.Source{
					Path:        "test.md",
					Content:     []byte("test"),
					ContentHash: "abcdef1234567890",
					Weight:      1.0,
				},
			},
		},
	}

	results, err := ext.Extract(context.Background(), snap, []config.SectionDecl{{Name: "Info", Budget: 1.0}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Content != "fenced response" {
		t.Errorf("content = %q, want %q", results[0].Content, "fenced response")
	}
}

func TestLLMExtractor_RequiresSections(t *testing.T) {
	ext := &LLMExtractor{Client: &mockClient{}}
	_, err := ext.Extract(context.Background(), protocol.Snapshot{}, nil)
	if err == nil {
		t.Fatal("expected error when no sections provided")
	}
}

func TestLLMExtractor_APIError(t *testing.T) {
	client := &mockClient{err: fmt.Errorf("API unavailable")}
	ext := &LLMExtractor{Client: client}
	snap := protocol.Snapshot{
		Items: []protocol.SnapshotItem{
			{Source: protocol.Source{Path: "test.md", Content: []byte("test"), ContentHash: "abc123"}},
		},
	}
	_, err := ext.Extract(context.Background(), snap, []config.SectionDecl{{Name: "Info", Budget: 1.0}})
	if err == nil {
		t.Fatal("expected error on API failure")
	}
}
