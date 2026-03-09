package extract

import (
	"context"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/model"
)

type PassthroughExtractor struct{}

func (p *PassthroughExtractor) Extract(_ context.Context, snap model.Snapshot, _ []config.SectionDecl) ([]model.Extraction, error) {
	var extractions []model.Extraction
	for _, item := range snap.Items {
		extractions = append(extractions, model.Extraction{
			Content:     string(item.Source.Content),
			Section:     "",
			TokenCount:  item.Source.TokenCount,
			Weight:      item.Source.Weight,
			SourcePath:  item.Source.Path,
			ContentHash: item.Source.ContentHash,
			DeclOrder:   item.Source.DeclOrder,
		})
	}
	return extractions, nil
}
