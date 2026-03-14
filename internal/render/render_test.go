package render

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

func TestRender_ValidJSON(t *testing.T) {
	ts := time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC)
	artifact := protocol.Artifact{
		Meta: protocol.Meta{
			Generator:    "cs v0.1.0",
			Mode:         "deterministic",
			SnapshotHash: "abc123",
			ConfigHash:   "def456",
			Timestamp:    ts,
			Budget:       5000,
			TokensUsed:   3000,
		},
		Sections: []protocol.ArtifactSection{
			{Name: "Overview", Items: []protocol.Extraction{
				{Content: "test content", SourcePath: "a.md", Weight: 1.0, ContentHash: "abcdef1234567890", TokenCount: 10},
			}},
		},
		Omissions: []protocol.Omission{},
	}

	output, err := Render(artifact)
	if err != nil {
		t.Fatal(err)
	}

	var parsed protocol.Artifact
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if parsed.Meta.Generator != "cs v0.1.0" {
		t.Errorf("generator = %q, want %q", parsed.Meta.Generator, "cs v0.1.0")
	}
	if parsed.Meta.Mode != "deterministic" {
		t.Errorf("mode = %q, want %q", parsed.Meta.Mode, "deterministic")
	}
	if parsed.Meta.Budget != 5000 {
		t.Errorf("budget = %d, want 5000", parsed.Meta.Budget)
	}
	if parsed.Meta.TokensUsed != 3000 {
		t.Errorf("tokens_used = %d, want 3000", parsed.Meta.TokensUsed)
	}
}

func TestRender_SourcesInJSON(t *testing.T) {
	artifact := protocol.Artifact{
		Meta: protocol.Meta{
			Generator:    "cs v0.1.0",
			Mode:         "deterministic",
			SnapshotHash: "snap123",
			ConfigHash:   "cfg456",
			Timestamp:    time.Now(),
			Budget:       5000,
		},
		Sections: []protocol.ArtifactSection{
			{Name: "Info", Items: []protocol.Extraction{
				{Content: "hello world", SourcePath: "docs/a.md", Weight: 0.8, ContentHash: "abcdef1234567890", TokenCount: 5},
			}},
		},
		Omissions: []protocol.Omission{},
	}

	output, err := Render(artifact)
	if err != nil {
		t.Fatal(err)
	}

	var parsed protocol.Artifact
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(parsed.Sections) != 1 {
		t.Fatalf("sections = %d, want 1", len(parsed.Sections))
	}
	if parsed.Sections[0].Items[0].SourcePath != "docs/a.md" {
		t.Errorf("source = %q, want docs/a.md", parsed.Sections[0].Items[0].SourcePath)
	}
	if parsed.Sections[0].Items[0].Weight != 0.8 {
		t.Errorf("weight = %f, want 0.8", parsed.Sections[0].Items[0].Weight)
	}
}

func TestRender_OmissionsInJSON(t *testing.T) {
	artifact := protocol.Artifact{
		Meta: protocol.Meta{
			Generator:    "cs v0.1.0",
			Mode:         "deterministic",
			SnapshotHash: "snap123",
			ConfigHash:   "cfg456",
			Timestamp:    time.Now(),
			Budget:       5000,
		},
		Sections:  []protocol.ArtifactSection{},
		Omissions: []protocol.Omission{
			{SourcePath: "extra.md", Reason: "exceeded budget", Weight: 0.3, TokenCount: 500},
		},
	}

	output, err := Render(artifact)
	if err != nil {
		t.Fatal(err)
	}

	var parsed protocol.Artifact
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(parsed.Omissions) != 1 {
		t.Fatalf("omissions = %d, want 1", len(parsed.Omissions))
	}
	if parsed.Omissions[0].SourcePath != "extra.md" {
		t.Errorf("omitted source = %q, want extra.md", parsed.Omissions[0].SourcePath)
	}
	if parsed.Omissions[0].Reason != "exceeded budget" {
		t.Errorf("omission reason = %q, want %q", parsed.Omissions[0].Reason, "exceeded budget")
	}
}
