package render

import (
	"strings"
	"testing"
	"time"

	"github.com/anthonymartinovic/context-synth/internal/model"
)

func TestRender_FrontMatter(t *testing.T) {
	ts := time.Date(2026, 3, 8, 12, 0, 0, 0, time.UTC)
	artifact := model.Artifact{
		FrontMatter: model.FrontMatter{
			Generator:  "cs v0.1.0",
			Mode:       "deterministic",
			Snapshot:   "abc123",
			Config:     "def456",
			Timestamp:  ts,
			Budget:     5000,
			TokensUsed: 3000,
		},
		Sections: []model.ArtifactSection{
			{Items: []model.Extraction{
				{Content: "test content", SourcePath: "a.md", Weight: 1.0, ContentHash: "abcdef1234567890"},
			}},
		},
	}

	output := Render(artifact)

	if !strings.Contains(output, "generator: cs v0.1.0") {
		t.Error("missing generator in front matter")
	}
	if !strings.Contains(output, "mode: deterministic") {
		t.Error("missing mode in front matter")
	}
	if !strings.Contains(output, "budget: 5000") {
		t.Error("missing budget in front matter")
	}
	if !strings.Contains(output, "tokens_used: 3000") {
		t.Error("missing tokens_used in front matter")
	}
	if !strings.Contains(output, "---") {
		t.Error("missing front matter delimiters")
	}
}

func TestRender_SourceAnnotations(t *testing.T) {
	artifact := model.Artifact{
		FrontMatter: model.FrontMatter{
			Generator: "cs v0.1.0",
			Mode:      "deterministic",
			Snapshot:  "snap123",
			Config:    "cfg456",
			Timestamp: time.Now(),
			Budget:    5000,
		},
		Sections: []model.ArtifactSection{
			{Items: []model.Extraction{
				{Content: "hello world", SourcePath: "docs/a.md", Weight: 0.8, ContentHash: "abcdef1234567890"},
			}},
		},
	}

	output := Render(artifact)

	if !strings.Contains(output, "> source: docs/a.md | weight: 0.8 | hash: abcdef12") {
		t.Errorf("missing source annotation in output:\n%s", output)
	}
}

func TestRender_Omissions(t *testing.T) {
	artifact := model.Artifact{
		FrontMatter: model.FrontMatter{
			Generator: "cs v0.1.0",
			Mode:      "deterministic",
			Snapshot:  "snap123",
			Config:    "cfg456",
			Timestamp: time.Now(),
			Budget:    5000,
		},
		Omissions: []model.Omission{
			{SourcePath: "extra.md", Reason: "exceeded budget", Weight: 0.3, TokenCount: 500},
		},
	}

	output := Render(artifact)

	if !strings.Contains(output, "# Omissions") {
		t.Error("missing omissions section")
	}
	if !strings.Contains(output, "extra.md") {
		t.Error("missing omitted source path")
	}
	if !strings.Contains(output, "exceeded budget") {
		t.Error("missing omission reason")
	}
}
