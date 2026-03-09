package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/config"
)

func setupTestFixture(t *testing.T) (config.Config, Options) {
	t.Helper()
	dir := t.TempDir()

	highContent := "# High Priority\n\nCritical domain knowledge that must always appear in the context artifact. This content has the highest weight."
	lowContent := "# Low Priority\n\nAdditional background and supplementary information that is less critical."

	os.WriteFile(filepath.Join(dir, "high.md"), []byte(highContent), 0644)
	os.WriteFile(filepath.Join(dir, "low.md"), []byte(lowContent), 0644)

	cfgContent := `version: "1"
budget: 10000
sources:
  - path: high.md
    weight: 1.0
  - path: low.md
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

	return cfg, opts
}

func TestPipeline_AC1_HigherWeightFirst(t *testing.T) {
	cfg, opts := setupTestFixture(t)

	result, err := Run(context.Background(), cfg, opts)
	if err != nil {
		t.Fatal(err)
	}

	highIdx := strings.Index(result.Output, "Critical domain knowledge")
	lowIdx := strings.Index(result.Output, "Additional background")

	if highIdx == -1 {
		t.Fatal("high-weight content not in output")
	}
	if lowIdx == -1 {
		t.Fatal("low-weight content not in output")
	}
	if highIdx >= lowIdx {
		t.Error("AC1: higher-weight content should appear before lower-weight content")
	}
}

func TestPipeline_AC2_NeverExceedsBudget(t *testing.T) {
	cfg, opts := setupTestFixture(t)

	result, err := Run(context.Background(), cfg, opts)
	if err != nil {
		t.Fatal(err)
	}

	if result.Artifact.FrontMatter.TokensUsed > cfg.Budget {
		t.Errorf("AC2: tokens used %d exceeds budget %d",
			result.Artifact.FrontMatter.TokensUsed, cfg.Budget)
	}
}

func TestPipeline_AC3_ProvenanceTracking(t *testing.T) {
	cfg, opts := setupTestFixture(t)

	result, err := Run(context.Background(), cfg, opts)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result.Output, "> source:") {
		t.Error("AC3: output missing source annotations")
	}
	if !strings.Contains(result.Output, "| hash:") {
		t.Error("AC3: output missing content hashes")
	}
}

func TestPipeline_AC4_OmissionsRecorded(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "a.md"), []byte(strings.Repeat("word ", 100)), 0644)
	os.WriteFile(filepath.Join(dir, "b.md"), []byte(strings.Repeat("more ", 100)), 0644)

	cfgContent := `version: "1"
budget: 50
sources:
  - path: a.md
    weight: 1.0
  - path: b.md
    weight: 0.5
`
	cfgPath := filepath.Join(dir, "contextsynth.yml")
	os.WriteFile(cfgPath, []byte(cfgContent), 0644)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Run(context.Background(), cfg, Options{ConfigPath: cfgPath, NoLLM: true})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Artifact.Omissions) == 0 {
		t.Error("AC4: expected omissions when budget is tight")
	}
	if !strings.Contains(result.Output, "# Omissions") {
		t.Error("AC4: omissions section missing from output")
	}
	if !strings.Contains(result.Output, "exceeded budget") {
		t.Error("AC4: omission reason missing")
	}
}

func TestPipeline_AC5_NoLLMProducesArtifact(t *testing.T) {
	cfg, opts := setupTestFixture(t)
	opts.NoLLM = true

	result, err := Run(context.Background(), cfg, opts)
	if err != nil {
		t.Fatal(err)
	}

	if result.Output == "" {
		t.Error("AC5: --no-llm produced empty output")
	}
	if !strings.Contains(result.Output, "mode: deterministic") {
		t.Error("AC5: mode should be deterministic")
	}
	if !strings.Contains(result.Output, "snapshot:") {
		t.Error("AC5: missing snapshot hash")
	}
}
