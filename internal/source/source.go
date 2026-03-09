package source

import (
	"context"

	"github.com/anthonymartinovic/context-synth/internal/config"
)

type ResolvedSource struct {
	Path        string
	Content     []byte
	ContentHash string
	TokenCount  int
	Weight      float64
	DeclOrder   int
}

type Provider interface {
	Resolve(ctx context.Context, decl config.SourceDecl, declOrder int) ([]ResolvedSource, error)
}
