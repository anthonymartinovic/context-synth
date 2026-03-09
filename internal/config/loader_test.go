package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestConfig(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "contextsynth.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadValid(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs"), 0755)
	os.WriteFile(filepath.Join(dir, "docs", "a.md"), []byte("test"), 0644)

	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources:
  - path: docs/a.md
    weight: 0.9
sections:
  - name: Overview
    budget: 0.6
  - name: Details
    budget: 0.4
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("version = %q, want \"1\"", cfg.Version)
	}
	if cfg.Budget != 5000 {
		t.Errorf("budget = %d, want 5000", cfg.Budget)
	}
	if cfg.Output.Path != "Contextfile" {
		t.Errorf("output path = %q, want \"Contextfile\"", cfg.Output.Path)
	}
	if len(cfg.Sources) != 1 {
		t.Fatalf("sources len = %d, want 1", len(cfg.Sources))
	}
	expectedPath := filepath.Join(dir, "docs", "a.md")
	if cfg.Sources[0].Path != expectedPath {
		t.Errorf("source path = %q, want %q", cfg.Sources[0].Path, expectedPath)
	}
	if cfg.Sections[0].Budget != 0.6 {
		t.Errorf("section 0 budget = %f, want 0.6", cfg.Sections[0].Budget)
	}
}

func TestLoadMissingSources(t *testing.T) {
	dir := t.TempDir()
	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources: []
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestLoadBadWeight(t *testing.T) {
	dir := t.TempDir()
	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources:
  - path: docs/a.md
    weight: 1.5
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for weight > 1.0")
	}
}

func TestLoadMissingBudget(t *testing.T) {
	dir := t.TempDir()
	path := writeTestConfig(t, dir, `
version: "1"
sources:
  - path: docs/a.md
    weight: 0.5
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing budget")
	}
}

func TestLoadBadVersion(t *testing.T) {
	dir := t.TempDir()
	path := writeTestConfig(t, dir, `
version: "2"
budget: 5000
sources:
  - path: docs/a.md
    weight: 0.5
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestSectionBudgetNormalization(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs"), 0755)
	os.WriteFile(filepath.Join(dir, "docs", "a.md"), []byte("test"), 0644)

	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources:
  - path: docs/a.md
    weight: 1.0
sections:
  - name: A
    budget: 3
  - name: B
    budget: 7
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Sections[0].Budget != 0.3 {
		t.Errorf("section A budget = %f, want 0.3", cfg.Sections[0].Budget)
	}
	if cfg.Sections[1].Budget != 0.7 {
		t.Errorf("section B budget = %f, want 0.7", cfg.Sections[1].Budget)
	}
}
