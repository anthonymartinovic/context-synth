package rank

import (
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/model"
)

func TestRank_WeightOrder(t *testing.T) {
	extractions := []model.Extraction{
		{SourcePath: "low.md", Weight: 0.3, DeclOrder: 0},
		{SourcePath: "high.md", Weight: 0.9, DeclOrder: 1},
		{SourcePath: "mid.md", Weight: 0.5, DeclOrder: 2},
	}

	ranked := Rank(extractions)

	if ranked[0].SourcePath != "high.md" {
		t.Errorf("first = %q, want high.md", ranked[0].SourcePath)
	}
	if ranked[1].SourcePath != "mid.md" {
		t.Errorf("second = %q, want mid.md", ranked[1].SourcePath)
	}
	if ranked[2].SourcePath != "low.md" {
		t.Errorf("third = %q, want low.md", ranked[2].SourcePath)
	}
}

func TestRank_TiebreakByDeclOrder(t *testing.T) {
	extractions := []model.Extraction{
		{SourcePath: "second.md", Weight: 0.5, DeclOrder: 1},
		{SourcePath: "first.md", Weight: 0.5, DeclOrder: 0},
	}

	ranked := Rank(extractions)

	if ranked[0].SourcePath != "first.md" {
		t.Errorf("first = %q, want first.md (lower decl order)", ranked[0].SourcePath)
	}
	if ranked[1].SourcePath != "second.md" {
		t.Errorf("second = %q, want second.md", ranked[1].SourcePath)
	}
}

func TestRank_DoesNotMutateOriginal(t *testing.T) {
	extractions := []model.Extraction{
		{SourcePath: "b.md", Weight: 0.3, DeclOrder: 1},
		{SourcePath: "a.md", Weight: 0.9, DeclOrder: 0},
	}

	_ = Rank(extractions)

	if extractions[0].SourcePath != "b.md" {
		t.Error("original slice was mutated")
	}
}
