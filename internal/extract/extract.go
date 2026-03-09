package extract

import (
	"context"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/model"
)

type Extractor interface {
	Extract(ctx context.Context, snap model.Snapshot, sections []config.SectionDecl) ([]model.Extraction, error)
}
