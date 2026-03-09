package assemble

import (
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/model"
)

func TestAssemble_WithinBudget(t *testing.T) {
	ranked := []model.Extraction{
		{SourcePath: "a.md", TokenCount: 100, Weight: 1.0},
		{SourcePath: "b.md", TokenCount: 100, Weight: 0.5},
	}

	plan := Assemble(ranked, 200)
	if len(plan.Included) != 2 {
		t.Errorf("included = %d, want 2", len(plan.Included))
	}
	if len(plan.Omitted) != 0 {
		t.Errorf("omitted = %d, want 0", len(plan.Omitted))
	}
}

func TestAssemble_ExceedsBudget(t *testing.T) {
	ranked := []model.Extraction{
		{SourcePath: "a.md", TokenCount: 100, Weight: 1.0},
		{SourcePath: "b.md", TokenCount: 150, Weight: 0.5},
	}

	plan := Assemble(ranked, 200)
	if len(plan.Included) != 1 {
		t.Errorf("included = %d, want 1", len(plan.Included))
	}
	if plan.Included[0].SourcePath != "a.md" {
		t.Errorf("included source = %q, want a.md", plan.Included[0].SourcePath)
	}
	if len(plan.Omitted) != 1 {
		t.Fatalf("omitted = %d, want 1", len(plan.Omitted))
	}
	if plan.Omitted[0].Reason != "exceeded budget" {
		t.Errorf("omission reason = %q, want %q", plan.Omitted[0].Reason, "exceeded budget")
	}
}

func TestAssemble_NeverExceedsBudget(t *testing.T) {
	ranked := []model.Extraction{
		{SourcePath: "a.md", TokenCount: 500, Weight: 1.0},
		{SourcePath: "b.md", TokenCount: 400, Weight: 0.9},
		{SourcePath: "c.md", TokenCount: 300, Weight: 0.8},
	}

	plan := Assemble(ranked, 800)
	var totalTokens int
	for _, item := range plan.Included {
		totalTokens += item.TokenCount
	}
	if totalTokens > 800 {
		t.Errorf("total tokens %d exceeds budget 800", totalTokens)
	}
}
