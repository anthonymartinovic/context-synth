package pipeline

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
	"github.com/anthonymartinovic/context-synth/integrations/source/markdown"
)

func setupTestFixture(t *testing.T) (config.Config, Options, Deps) {
	t.Helper()
	dir := t.TempDir()

	highContent := "# High Priority\n\nCritical domain knowledge that must always appear in the context artifact. This content has the highest weight."
	lowContent := "# Low Priority\n\nAdditional background and supplementary information that is less critical."

	os.WriteFile(filepath.Join(dir, "high.md"), []byte(highContent), 0644)
	os.WriteFile(filepath.Join(dir, "low.md"), []byte(lowContent), 0644)

	cfgContent := `version: "1"
budget: 10000
sources:
  - markdown: high.md
    weight: 1.0
  - markdown: low.md
    weight: 0.3
`
	cfgPath := filepath.Join(dir, "contextsynth.yml")
	os.WriteFile(cfgPath, []byte(cfgContent), 0644)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		ConfigPath: cfgPath,
		NoLLM:      true,
	}

	deps := Deps{
		Provider: &markdown.Provider{},
	}

	return cfg, opts, deps
}

func TestPipeline_AC1_HigherWeightFirst(t *testing.T) {
	cfg, opts, deps := setupTestFixture(t)

	result, err := Run(context.Background(), cfg, opts, deps)
	if err != nil {
		t.Fatal(err)
	}

	var parsed protocol.Artifact
	if err := json.Unmarshal([]byte(result.Output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(parsed.Sections) == 0 || len(parsed.Sections[0].Items) < 2 {
		t.Fatal("expected at least 2 items in output")
	}

	first := parsed.Sections[0].Items[0]
	second := parsed.Sections[0].Items[1]

	if first.Weight < second.Weight {
		t.Error("AC1: higher-weight content should appear before lower-weight content")
	}
}

func TestPipeline_AC2_NeverExceedsBudget(t *testing.T) {
	cfg, opts, deps := setupTestFixture(t)

	result, err := Run(context.Background(), cfg, opts, deps)
	if err != nil {
		t.Fatal(err)
	}

	if result.Artifact.Meta.TokensUsed > cfg.Budget {
		t.Errorf("AC2: tokens used %d exceeds budget %d",
			result.Artifact.Meta.TokensUsed, cfg.Budget)
	}
}

func TestPipeline_AC3_ProvenanceTracking(t *testing.T) {
	cfg, opts, deps := setupTestFixture(t)

	result, err := Run(context.Background(), cfg, opts, deps)
	if err != nil {
		t.Fatal(err)
	}

	var parsed protocol.Artifact
	if err := json.Unmarshal([]byte(result.Output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	for _, sec := range parsed.Sections {
		for _, item := range sec.Items {
			if item.SourcePath == "" {
				t.Error("AC3: item missing source path")
			}
			if item.ContentHash == "" {
				t.Error("AC3: item missing content hash")
			}
		}
	}
}

func TestPipeline_AC4_OmissionsRecorded(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "a.md"), []byte(repeatWords("word ", 100)), 0644)
	os.WriteFile(filepath.Join(dir, "b.md"), []byte(repeatWords("more ", 100)), 0644)

	cfgContent := `version: "1"
budget: 50
sources:
  - markdown: a.md
    weight: 1.0
  - markdown: b.md
    weight: 0.5
`
	cfgPath := filepath.Join(dir, "contextsynth.yml")
	os.WriteFile(cfgPath, []byte(cfgContent), 0644)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	deps := Deps{Provider: &markdown.Provider{}}
	result, err := Run(context.Background(), cfg, Options{ConfigPath: cfgPath, NoLLM: true}, deps)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Artifact.Omissions) == 0 {
		t.Error("AC4: expected omissions when budget is tight")
	}

	var parsed protocol.Artifact
	if err := json.Unmarshal([]byte(result.Output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	foundOmission := false
	for _, o := range parsed.Omissions {
		if o.Reason == "exceeded budget" {
			foundOmission = true
			break
		}
	}
	if !foundOmission {
		t.Error("AC4: omission reason 'exceeded budget' not found in JSON output")
	}
}

func TestPipeline_AC5_NoLLMProducesArtifact(t *testing.T) {
	cfg, opts, deps := setupTestFixture(t)
	opts.NoLLM = true

	result, err := Run(context.Background(), cfg, opts, deps)
	if err != nil {
		t.Fatal(err)
	}

	if result.Output == "" {
		t.Error("AC5: --no-llm produced empty output")
	}

	var parsed protocol.Artifact
	if err := json.Unmarshal([]byte(result.Output), &parsed); err != nil {
		t.Fatalf("AC5: output is not valid JSON: %v", err)
	}

	if parsed.Meta.Mode != "deterministic" {
		t.Errorf("AC5: mode = %q, want deterministic", parsed.Meta.Mode)
	}
	if parsed.Meta.SnapshotHash == "" {
		t.Error("AC5: missing snapshot hash")
	}
}

func repeatWords(word string, n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += word
	}
	return s
}
