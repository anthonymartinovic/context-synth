package extract

import (
	"context"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type Extractor interface {
	Extract(ctx context.Context, snap protocol.Snapshot, sections []config.SectionDecl) ([]protocol.Extraction, error)
}
