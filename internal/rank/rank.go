package rank

import (
	"sort"

	"github.com/anthonymartinovic/context-synth/internal/model"
)

func Rank(extractions []model.Extraction) []model.Extraction {
	ranked := make([]model.Extraction, len(extractions))
	copy(ranked, extractions)

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Weight != ranked[j].Weight {
			return ranked[i].Weight > ranked[j].Weight
		}
		return ranked[i].DeclOrder < ranked[j].DeclOrder
	})

	return ranked
}
