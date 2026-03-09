package snapshot

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/model"
	"github.com/anthonymartinovic/context-synth/internal/source"
)

func Build(ctx context.Context, cfg config.Config, provider source.Provider) (model.Snapshot, error) {
	var allItems []model.SnapshotItem

	for i, decl := range cfg.Sources {
		resolved, err := provider.Resolve(ctx, decl, i)
		if err != nil {
			return model.Snapshot{}, fmt.Errorf("resolving source %d (%s): %w", i, decl.Path, err)
		}
		for _, r := range resolved {
			allItems = append(allItems, model.SnapshotItem{
				Source: model.Source{
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
	return model.Snapshot{Items: allItems, Hash: hash}, nil
}

func computeSnapshotHash(items []model.SnapshotItem) string {
	h := sha256.New()
	for _, item := range items {
		h.Write([]byte(item.Source.ContentHash))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
