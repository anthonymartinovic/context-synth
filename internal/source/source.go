package source

import (
	"context"
	"fmt"

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

type MultiProvider struct {
	Markdown Provider
	Audio    Provider
}

func (mp *MultiProvider) Resolve(ctx context.Context, decl config.SourceDecl, declOrder int) ([]ResolvedSource, error) {
	if decl.Audio != "" {
		if mp.Audio == nil {
			return nil, fmt.Errorf("audio source declared but no audio provider configured")
		}
		return mp.Audio.Resolve(ctx, decl, declOrder)
	}
	if mp.Markdown == nil {
		return nil, fmt.Errorf("markdown source declared but no markdown provider configured")
	}
	return mp.Markdown.Resolve(ctx, decl, declOrder)
}
