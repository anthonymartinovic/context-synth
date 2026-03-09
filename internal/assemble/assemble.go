package assemble

import (
	"github.com/anthonymartinovic/context-synth/internal/model"
)

func Assemble(ranked []model.Extraction, budget int) model.SynthPlan {
	var plan model.SynthPlan
	tokensUsed := 0

	for _, item := range ranked {
		if tokensUsed+item.TokenCount > budget {
			plan.Omitted = append(plan.Omitted, model.Omission{
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

func AssembleWithSections(ranked []model.Extraction, sections []model.SectionBudget, totalBudget int) model.SynthPlan {
	sectionTokens := make(map[string]int)
	sectionBudgets := make(map[string]int)
	for _, s := range sections {
		sectionBudgets[s.Name] = int(s.Proportion * float64(totalBudget))
	}

	var plan model.SynthPlan
	for _, item := range ranked {
		budget, hasBudget := sectionBudgets[item.Section]
		if !hasBudget {
			plan.Omitted = append(plan.Omitted, model.Omission{
				SourcePath: item.SourcePath,
				Reason:     "no matching section",
				Weight:     item.Weight,
				TokenCount: item.TokenCount,
			})
			continue
		}
		if sectionTokens[item.Section]+item.TokenCount > budget {
			plan.Omitted = append(plan.Omitted, model.Omission{
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
