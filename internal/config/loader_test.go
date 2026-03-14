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
  - markdown: docs/a.md
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
	if cfg.Output.Path != "artifact.json" {
		t.Errorf("output path = %q, want \"artifact.json\"", cfg.Output.Path)
	}
	if len(cfg.Sources) != 1 {
		t.Fatalf("sources len = %d, want 1", len(cfg.Sources))
	}
	expectedMarkdown := filepath.Join(dir, "docs", "a.md")
	if cfg.Sources[0].Markdown != expectedMarkdown {
		t.Errorf("source markdown = %q, want %q", cfg.Sources[0].Markdown, expectedMarkdown)
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
  - markdown: docs/a.md
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
  - markdown: docs/a.md
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
  - markdown: docs/a.md
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
  - markdown: docs/a.md
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

func TestLoadAudioSource(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "tracks"), 0755)
	os.WriteFile(filepath.Join(dir, "tracks", "song.mp3"), []byte("fake"), 0644)
	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources:
  - audio: tracks/song.mp3
    weight: 0.8
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedAudio := filepath.Join(dir, "tracks", "song.mp3")
	if cfg.Sources[0].Audio != expectedAudio {
		t.Errorf("audio = %q, want %q", cfg.Sources[0].Audio, expectedAudio)
	}
}

func TestLoadMixedSourcesWithAudio(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "vibe.md"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(dir, "track.mp3"), []byte("fake"), 0644)
	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources:
  - markdown: vibe.md
    weight: 1.0
  - audio: track.mp3
    weight: 0.7
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 2 {
		t.Fatalf("sources = %d, want 2", len(cfg.Sources))
	}
	if cfg.Sources[1].Audio == "" {
		t.Error("expected audio field to be set on source 1")
	}
}

func TestLoadMarkdownAndAudioConflict(t *testing.T) {
	dir := t.TempDir()
	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources:
  - markdown: docs/a.md
    audio: track.mp3
    weight: 0.8
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for source with both markdown and audio")
	}
}

func TestLoadLlamaConfig(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.md"), []byte("test"), 0644)
	path := writeTestConfig(t, dir, `
version: "1"
budget: 5000
sources:
  - markdown: a.md
    weight: 1.0
llm:
  provider: llama
  model: llama3.1:8b
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LLM.Provider != "llama" {
		t.Errorf("provider = %q, want llama", cfg.LLM.Provider)
	}
	if cfg.LLM.Model != "llama3.1:8b" {
		t.Errorf("model = %q, want llama3.1:8b", cfg.LLM.Model)
	}
}
