package verify

import "github.com/anthonymartinovic/context-synth/internal/protocol"

type VerificationData struct {
	Mode     string
	Included []IncludedItem
	Omitted  []protocol.Omission
}

type IncludedItem struct {
	SourcePath  string
	ContentHash string
	Weight      float64
	TokenCount  int
	Section     string
}

func Verify(plan protocol.SynthPlan, mode string) VerificationData {
	var included []IncludedItem
	for _, item := range plan.Included {
		included = append(included, IncludedItem{
			SourcePath:  item.SourcePath,
			ContentHash: item.ContentHash,
			Weight:      item.Weight,
			TokenCount:  item.TokenCount,
			Section:     item.Section,
		})
	}
	return VerificationData{
		Mode:     mode,
		Included: included,
		Omitted:  plan.Omitted,
	}
}
