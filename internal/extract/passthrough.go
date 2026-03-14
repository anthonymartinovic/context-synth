package extract

import (
	"context"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type PassthroughExtractor struct{}

func (p *PassthroughExtractor) Extract(_ context.Context, snap protocol.Snapshot, _ []config.SectionDecl) ([]protocol.Extraction, error) {
	var extractions []protocol.Extraction
	for _, item := range snap.Items {
		extractions = append(extractions, protocol.Extraction{
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
