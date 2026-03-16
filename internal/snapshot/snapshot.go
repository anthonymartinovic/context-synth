package snapshot

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
	"github.com/anthonymartinovic/context-synth/internal/source"
)

func Build(ctx context.Context, cfg config.Config, provider source.Provider) (protocol.Snapshot, error) {
	var allItems []protocol.SnapshotItem

	for i, decl := range cfg.Sources {
		resolved, err := provider.Resolve(ctx, decl, i)
		if err != nil {
			label := decl.Markdown
			if label == "" {
				label = decl.Audio
			}
			return protocol.Snapshot{}, fmt.Errorf("resolving source %d (%s): %w", i, label, err)
		}
		for _, r := range resolved {
			allItems = append(allItems, protocol.SnapshotItem{
				Source: protocol.Source{
					Path:        r.Path,
					Content:     r.Content,
					ContentHash: r.ContentHash,
					TokenCount:  r.TokenCount,
					Weight:      r.Weight,
					DeclOrder:   r.DeclOrder,
				},
				Provenance: r.Path,
			})
		}
	}

	sort.SliceStable(allItems, func(i, j int) bool {
		return allItems[i].Source.DeclOrder < allItems[j].Source.DeclOrder
	})

	hash := computeSnapshotHash(allItems)
	return protocol.Snapshot{Items: allItems, Hash: hash}, nil
}

func computeSnapshotHash(items []protocol.SnapshotItem) string {
	h := sha256.New()
	for _, item := range items {
		h.Write([]byte(item.Source.ContentHash))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
