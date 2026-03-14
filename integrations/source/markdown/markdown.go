package markdown

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/source"
	"github.com/anthonymartinovic/context-synth/internal/token"
)

type Provider struct{}

func (m *Provider) Resolve(_ context.Context, decl config.SourceDecl, declOrder int) ([]source.ResolvedSource, error) {
	matches, err := filepath.Glob(decl.Markdown)
	if err != nil {
		return nil, fmt.Errorf("glob %q: %w", decl.Markdown, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no files matched %q", decl.Markdown)
	}

	sort.Strings(matches)

	var sources []source.ResolvedSource
	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("stat %q: %w", path, err)
		}
		if info.IsDir() {
			continue
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}
		hash := sha256.Sum256(content)
		sources = append(sources, source.ResolvedSource{
			Path:        path,
			Content:     content,
			ContentHash: fmt.Sprintf("%x", hash),
			TokenCount:  token.Estimate(string(content)),
			Weight:      decl.Weight,
			DeclOrder:   declOrder,
		})
	}
	return sources, nil
}
