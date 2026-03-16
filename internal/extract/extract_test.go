package extract

import (
	"context"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

func TestPassthroughExtractor(t *testing.T) {
	snap := protocol.Snapshot{
		Items: []protocol.SnapshotItem{
			{
				Source: protocol.Source{
					Path:        "a.md",
					Content:     []byte("alpha content"),
					ContentHash: "abc123",
					TokenCount:  10,
					Weight:      1.0,
					DeclOrder:   0,
				},
				Provenance: "a.md",
			},
			{
				Source: protocol.Source{
					Path:        "b.md",
					Content:     []byte("bravo content"),
					ContentHash: "def456",
					TokenCount:  8,
					Weight:      0.5,
					DeclOrder:   1,
				},
				Provenance: "b.md",
			},
		},
	}

	ext := &PassthroughExtractor{}
	results, err := ext.Extract(context.Background(), snap, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d extractions, want 2", len(results))
	}
	if results[0].Content != "alpha content" {
		t.Errorf("first content = %q, want %q", results[0].Content, "alpha content")
	}
	if results[0].Weight != 1.0 {
		t.Errorf("first weight = %f, want 1.0", results[0].Weight)
	}
	if results[0].SourcePath != "a.md" {
		t.Errorf("first source = %q, want %q", results[0].SourcePath, "a.md")
	}
	if results[0].ContentHash != "abc123" {
		t.Errorf("first hash = %q, want %q", results[0].ContentHash, "abc123")
	}
	if results[0].Section != "" {
		t.Errorf("section should be empty in passthrough, got %q", results[0].Section)
	}
}
