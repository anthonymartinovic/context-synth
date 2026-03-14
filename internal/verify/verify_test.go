package verify

import (
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

func TestVerify_IncludedAndOmitted(t *testing.T) {
	plan := protocol.SynthPlan{
		Included: []protocol.Extraction{
			{SourcePath: "a.md", ContentHash: "abc123", Weight: 1.0, TokenCount: 100},
		},
		Omitted: []protocol.Omission{
			{SourcePath: "b.md", Reason: "exceeded budget", Weight: 0.5, TokenCount: 200},
		},
	}

	data := Verify(plan, "deterministic")

	if data.Mode != "deterministic" {
		t.Errorf("mode = %q, want deterministic", data.Mode)
	}
	if len(data.Included) != 1 {
		t.Fatalf("included = %d, want 1", len(data.Included))
	}
	if data.Included[0].SourcePath != "a.md" {
		t.Errorf("included source = %q, want a.md", data.Included[0].SourcePath)
	}
	if len(data.Omitted) != 1 {
		t.Fatalf("omitted = %d, want 1", len(data.Omitted))
	}
	if data.Omitted[0].Reason != "exceeded budget" {
		t.Errorf("omission reason = %q", data.Omitted[0].Reason)
	}
}
