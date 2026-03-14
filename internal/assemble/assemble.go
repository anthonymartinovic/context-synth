package assemble

import (
	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

func Assemble(ranked []protocol.Extraction, budget int) protocol.SynthPlan {
	var plan protocol.SynthPlan
	tokensUsed := 0

	for _, item := range ranked {
		if tokensUsed+item.TokenCount > budget {
			plan.Omitted = append(plan.Omitted, protocol.Omission{
				SourcePath: item.SourcePath,
				Reason:     "exceeded budget",
				Weight:     item.Weight,
				TokenCount: item.TokenCount,
			})
			continue
		}
		plan.Included = append(plan.Included, item)
		tokensUsed += item.TokenCount
	}

	return plan
}

func AssembleWithSections(ranked []protocol.Extraction, sections []protocol.SectionBudget, totalBudget int) protocol.SynthPlan {
	sectionTokens := make(map[string]int)
	sectionBudgets := make(map[string]int)
	for _, s := range sections {
		sectionBudgets[s.Name] = int(s.Proportion * float64(totalBudget))
	}

	var plan protocol.SynthPlan
	for _, item := range ranked {
		budget, hasBudget := sectionBudgets[item.Section]
		if !hasBudget {
			plan.Omitted = append(plan.Omitted, protocol.Omission{
				SourcePath: item.SourcePath,
				Reason:     "no matching section",
				Weight:     item.Weight,
				TokenCount: item.TokenCount,
			})
			continue
		}
		if sectionTokens[item.Section]+item.TokenCount > budget {
			plan.Omitted = append(plan.Omitted, protocol.Omission{
				SourcePath: item.SourcePath,
				Reason:     "exceeded section budget",
				Weight:     item.Weight,
				TokenCount: item.TokenCount,
			})
			continue
		}
		plan.Included = append(plan.Included, item)
		sectionTokens[item.Section] += item.TokenCount
	}

	return plan
}
