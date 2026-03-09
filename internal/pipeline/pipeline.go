package pipeline

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/anthonymartinovic/context-synth/internal/assemble"
	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/extract"
	"github.com/anthonymartinovic/context-synth/internal/llm"
	"github.com/anthonymartinovic/context-synth/internal/model"
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
	Artifact model.Artifact
	Output   string
}

func Run(ctx context.Context, cfg config.Config, opts Options) (Result, error) {
	provider := &source.MarkdownProvider{}
	snap, err := snapshot.Build(ctx, cfg, provider)
	if err != nil {
		return Result{}, fmt.Errorf("snapshot: %w", err)
	}

	var extractor extract.Extractor
	mode := "deterministic"

	if cfg.LLM != nil && !opts.NoLLM {
		mode = "full"
		llmExtractor, err := buildLLMExtractor(cfg)
		if err != nil {
			return Result{}, fmt.Errorf("llm extractor: %w", err)
		}
		extractor = llmExtractor
	} else {
		extractor = &extract.PassthroughExtractor{}
	}

	extractions, err := extractor.Extract(ctx, snap, cfg.Sections)
	if err != nil {
		return Result{}, fmt.Errorf("extract: %w", err)
	}

	ranked := rank.Rank(extractions)

	var plan model.SynthPlan
	if mode == "full" && len(cfg.Sections) > 0 {
		var sectionBudgets []model.SectionBudget
		for _, s := range cfg.Sections {
			sectionBudgets = append(sectionBudgets, model.SectionBudget{
				Name:       s.Name,
				Proportion: s.Budget,
			})
		}
		plan = assemble.AssembleWithSections(ranked, sectionBudgets, cfg.Budget)
	} else {
		plan = assemble.Assemble(ranked, cfg.Budget)
	}

	_ = verify.Verify(plan, mode)

	configHash := computeConfigHash(opts.ConfigPath)

	var tokensUsed int
	for _, item := range plan.Included {
		tokensUsed += item.TokenCount
	}

	var sections []model.ArtifactSection
	if mode == "full" && len(cfg.Sections) > 0 {
		sectionMap := make(map[string]*model.ArtifactSection)
		for _, s := range cfg.Sections {
			sec := model.ArtifactSection{Name: s.Name}
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
		sections = []model.ArtifactSection{
			{Name: "", Items: plan.Included},
		}
	}

	artifact := model.Artifact{
		FrontMatter: model.FrontMatter{
			Generator:  "cs v0.1.0",
			Mode:       mode,
			Snapshot:   snap.Hash,
			Config:     configHash,
			Timestamp:  time.Now(),
			Budget:     cfg.Budget,
			TokensUsed: tokensUsed,
		},
		Sections:  sections,
		Omissions: plan.Omitted,
	}

	output := render.Render(artifact)

	return Result{Artifact: artifact, Output: output}, nil
}

func computeConfigHash(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func buildLLMExtractor(cfg config.Config) (extract.Extractor, error) {
	client, err := llm.NewGeminiClient(cfg.LLM.Model, cfg.LLM.APIKeyEnv)
	if err != nil {
		return nil, err
	}
	return &extract.LLMExtractor{Client: client}, nil
}
