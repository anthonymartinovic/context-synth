package snapshot

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/source"
)

func TestBuild_Deterministic(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.md"), []byte("alpha content"), 0644)
	os.WriteFile(filepath.Join(dir, "b.md"), []byte("bravo content"), 0644)

	cfg := config.Config{
		Sources: []config.SourceDecl{
			{Path: filepath.Join(dir, "a.md"), Weight: 1.0},
			{Path: filepath.Join(dir, "b.md"), Weight: 0.5},
		},
	}

	provider := &source.MarkdownProvider{}
	snap1, err := Build(context.Background(), cfg, provider)
	if err != nil {
		t.Fatal(err)
	}
	snap2, err := Build(context.Background(), cfg, provider)
	if err != nil {
		t.Fatal(err)
	}

	if snap1.Hash != snap2.Hash {
		t.Errorf("snapshot hashes differ: %s vs %s", snap1.Hash, snap2.Hash)
	}
	if len(snap1.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(snap1.Items))
	}
}

func TestBuild_DeclOrder(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "first.md"), []byte("first"), 0644)
	os.WriteFile(filepath.Join(dir, "second.md"), []byte("second"), 0644)

	cfg := config.Config{
		Sources: []config.SourceDecl{
			{Path: filepath.Join(dir, "first.md"), Weight: 0.3},
			{Path: filepath.Join(dir, "second.md"), Weight: 0.9},
		},
	}

	provider := &source.MarkdownProvider{}
	snap, err := Build(context.Background(), cfg, provider)
	if err != nil {
		t.Fatal(err)
	}

	if snap.Items[0].Source.DeclOrder != 0 {
		t.Errorf("first item decl order = %d, want 0", snap.Items[0].Source.DeclOrder)
	}
	if snap.Items[1].Source.DeclOrder != 1 {
		t.Errorf("second item decl order = %d, want 1", snap.Items[1].Source.DeclOrder)
	}
}
