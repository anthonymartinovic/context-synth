package pipeline

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/anthonymartinovic/context-synth/internal/assemble"
	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/extract"
	"github.com/anthonymartinovic/context-synth/internal/llm"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
	"github.com/anthonymartinovic/context-synth/internal/rank"
	"github.com/anthonymartinovic/context-synth/internal/render"
	"github.com/anthonymartinovic/context-synth/internal/snapshot"
	"github.com/anthonymartinovic/context-synth/internal/source"
	"github.com/anthonymartinovic/context-synth/internal/verify"
)

type Options struct {
	ConfigPath string
	OutputPath string
	NoLLM      bool
	DryRun     bool
	Verbose    bool
}

type Result struct {
	Artifact protocol.Artifact
	Output   string
}

type Deps struct {
	Provider  source.Provider
	LLMClient llm.Client
}

func Run(ctx context.Context, cfg config.Config, opts Options, deps Deps) (Result, error) {
	snap, err := snapshot.Build(ctx, cfg, deps.Provider)
	if err != nil {
		return Result{}, fmt.Errorf("snapshot: %w", err)
	}

	var extractor extract.Extractor
	mode := "deterministic"

	if deps.LLMClient != nil && cfg.LLM != nil && !opts.NoLLM {
		mode = "full"
		extractor = &extract.LLMExtractor{Client: deps.LLMClient}
	} else {
		extractor = &extract.PassthroughExtractor{}
	}

	extractions, err := extractor.Extract(ctx, snap, cfg.Sections)
	if err != nil {
		return Result{}, fmt.Errorf("extract: %w", err)
	}

	ranked := rank.Rank(extractions)

	var plan protocol.SynthPlan
	if mode == "full" && len(cfg.Sections) > 0 {
		var sectionBudgets []protocol.SectionBudget
		for _, s := range cfg.Sections {
			sectionBudgets = append(sectionBudgets, protocol.SectionBudget{
				Name:       s.Name,
				Proportion: s.Budget,
			})
		}
		plan = assemble.AssembleWithSections(ranked, sectionBudgets, cfg.Budget)
	} else {
		plan = assemble.Assemble(ranked, cfg.Budget)
	}

	vd := verify.Verify(plan, mode)
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "verify: mode=%s included=%d omitted=%d\n",
			vd.Mode, len(vd.Included), len(vd.Omitted))
	}

	configHash := computeConfigHash(opts.ConfigPath)

	var tokensUsed int
	for _, item := range plan.Included {
		tokensUsed += item.TokenCount
	}

	var sections []protocol.ArtifactSection
	if mode == "full" && len(cfg.Sections) > 0 {
		sectionMap := make(map[string]*protocol.ArtifactSection)
		for _, s := range cfg.Sections {
			sec := protocol.ArtifactSection{Name: s.Name}
			sectionMap[s.Name] = &sec
		}
		for _, item := range plan.Included {
			if sec, ok := sectionMap[item.Section]; ok {
				sec.Items = append(sec.Items, item)
			}
		}
		for _, s := range cfg.Sections {
			sec := sectionMap[s.Name]
			if len(sec.Items) > 0 {
				sections = append(sections, *sec)
			}
		}
	} else {
		sections = []protocol.ArtifactSection{
			{Name: "", Items: plan.Included},
		}
	}

	artifact := protocol.Artifact{
		Meta: protocol.Meta{
			Generator:    "cs " + binaryVersion(),
			Mode:         mode,
			SnapshotHash: snap.Hash,
			ConfigHash:   configHash,
			Timestamp:    time.Now(),
			Budget:       cfg.Budget,
			TokensUsed:   tokensUsed,
		},
		Sections:  sections,
		Omissions: plan.Omitted,
	}

	output, err := render.Render(artifact)
	if err != nil {
		return Result{}, fmt.Errorf("render: %w", err)
	}

	return Result{Artifact: artifact, Output: output}, nil
}

func binaryVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "(devel)"
}

func computeConfigHash(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
