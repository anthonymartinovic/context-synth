package source

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/config"
)

func TestMarkdownResolve_SingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	os.WriteFile(path, []byte("hello world"), 0644)

	p := &MarkdownProvider{}
	results, err := p.Resolve(context.Background(), config.SourceDecl{Path: path, Weight: 0.8}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Weight != 0.8 {
		t.Errorf("weight = %f, want 0.8", results[0].Weight)
	}
	if results[0].TokenCount <= 0 {
		t.Errorf("token count should be positive, got %d", results[0].TokenCount)
	}
	if results[0].ContentHash == "" {
		t.Error("content hash should not be empty")
	}
}

func TestMarkdownResolve_Glob(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.md"), []byte("alpha"), 0644)
	os.WriteFile(filepath.Join(dir, "b.md"), []byte("bravo"), 0644)
	os.WriteFile(filepath.Join(dir, "c.txt"), []byte("charlie"), 0644)

	p := &MarkdownProvider{}
	results, err := p.Resolve(context.Background(), config.SourceDecl{
		Path:   filepath.Join(dir, "*.md"),
		Weight: 0.5,
	}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
}

func TestMarkdownResolve_NoMatch(t *testing.T) {
	p := &MarkdownProvider{}
	_, err := p.Resolve(context.Background(), config.SourceDecl{
		Path:   "/nonexistent/*.md",
		Weight: 0.5,
	}, 0)
	if err == nil {
		t.Fatal("expected error for no matches")
	}
}

func TestMarkdownResolve_StableHash(t *testing.T) {
	dir := t.TempDir()
	content := []byte("consistent content")
	path := filepath.Join(dir, "stable.md")
	os.WriteFile(path, content, 0644)

	p := &MarkdownProvider{}
	r1, _ := p.Resolve(context.Background(), config.SourceDecl{Path: path, Weight: 1.0}, 0)
	r2, _ := p.Resolve(context.Background(), config.SourceDecl{Path: path, Weight: 1.0}, 0)

	if r1[0].ContentHash != r2[0].ContentHash {
		t.Errorf("hashes differ for same content: %s vs %s", r1[0].ContentHash, r2[0].ContentHash)
	}
}
