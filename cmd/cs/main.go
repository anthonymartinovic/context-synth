package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"path/filepath"

	"github.com/anthonymartinovic/context-synth/integrations/llm/llama"
	"github.com/anthonymartinovic/context-synth/integrations/source/audio"
	"github.com/anthonymartinovic/context-synth/integrations/source/markdown"
	"github.com/anthonymartinovic/context-synth/internal/adapter"
	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/pipeline"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
	"github.com/anthonymartinovic/context-synth/internal/runtime"
	"github.com/anthonymartinovic/context-synth/internal/snapshot"
	"github.com/anthonymartinovic/context-synth/internal/source"
	"github.com/spf13/cobra"
)

var cfgFile string

func main() {
	root := &cobra.Command{
		Use:   "cs",
		Short: "Context Synth — compile governed context artifacts and resolve capability graphs",
	}
	root.PersistentFlags().StringVar(&cfgFile, "config", "contextsynth.yml", "path to config file")

	root.AddCommand(snapCmd())
	root.AddCommand(synthCmd())
	root.AddCommand(runCmd())
	root.AddCommand(initCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func buildProviders(cfg config.Config) source.MultiProvider {
	mp := source.MultiProvider{
		Markdown: &markdown.Provider{},
	}

	for _, s := range cfg.Sources {
		if s.Audio != "" {
			mp.Audio = &audio.Provider{
				AnalyzerPath: os.Getenv("CS_AUDIO_ANALYZER"),
			}
			break
		}
	}

	return mp
}

func snapCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "snap",
		Short: "Resolve sources and print the snapshot table",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			providers := buildProviders(cfg)
			snap, err := snapshot.Build(context.Background(), cfg, &providers)
			if err != nil {
				return err
			}

			var totalTokens int
			for _, item := range snap.Items {
				totalTokens += item.Source.TokenCount
			}

			fmt.Printf("%-40s %6s %8s %s\n", "Source", "Weight", "Tokens", "Hash")
			for _, item := range snap.Items {
				fmt.Printf("%-40s %6.2f %8s %s\n",
					item.Source.Path,
					item.Source.Weight,
					formatNumber(item.Source.TokenCount),
					item.Source.ContentHash[:8],
				)
			}
			fmt.Printf("\n%d sources | %s tokens estimated | budget: %s\n",
				len(snap.Items),
				formatNumber(totalTokens),
				formatNumber(cfg.Budget),
			)
			return nil
		},
	}
}

func synthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synth",
		Short: "Run the pipeline and write the context artifact (JSON)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}

			noLLM, _ := cmd.Flags().GetBool("no-llm")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			verbose, _ := cmd.Flags().GetBool("verbose")
			outputFlag, _ := cmd.Flags().GetString("output")

			outputPath := cfg.Output.Path
			if outputFlag != "" {
				outputPath = outputFlag
			}

			opts := pipeline.Options{
				ConfigPath: cfgFile,
				OutputPath: outputPath,
				NoLLM:      noLLM,
				DryRun:     dryRun,
				Verbose:    verbose,
			}

			providers := buildProviders(cfg)
			deps := pipeline.Deps{
				Provider: &providers,
			}

			if cfg.LLM != nil && cfg.LLM.Provider == "llama" && !noLLM {
				deps.LLMClient = llama.NewClient(cfg.LLM.Model, cfg.LLM.BaseURL)
			}

			result, err := pipeline.Run(context.Background(), cfg, opts, deps)
			if err != nil {
				return err
			}

			if verbose {
				meta := result.Artifact.Meta
				var includedCount int
				for _, sec := range result.Artifact.Sections {
					includedCount += len(sec.Items)
				}
				fmt.Fprintf(os.Stderr, "mode: %s | tokens: %d/%d | sources: %d included, %d omitted\n",
					meta.Mode, meta.TokensUsed, meta.Budget,
					includedCount, len(result.Artifact.Omissions))
			}

			if dryRun {
				fmt.Print(result.Output)
				return nil
			}

			if err := os.WriteFile(outputPath, []byte(result.Output), 0644); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}
			fmt.Fprintf(os.Stderr, "wrote %s\n", outputPath)
			return nil
		},
	}
	cmd.Flags().String("output", "", "output file path (overrides config)")
	cmd.Flags().Bool("no-llm", false, "deterministic mode without LLM")
	cmd.Flags().Bool("dry-run", false, "preview without writing")
	cmd.Flags().Bool("verbose", false, "verbose output")
	return cmd
}

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Load artifact, resolve capability graph, and project output",
		RunE: func(cmd *cobra.Command, args []string) error {
			artifactPath, _ := cmd.Flags().GetString("artifact")
			capsPath, _ := cmd.Flags().GetString("capabilities")
			adapterCmd, _ := cmd.Flags().GetString("adapter")
			verbose, _ := cmd.Flags().GetBool("verbose")

			artifactData, err := os.ReadFile(artifactPath)
			if err != nil {
				return fmt.Errorf("reading artifact: %w", err)
			}
			var artifact protocol.Artifact
			if err := json.Unmarshal(artifactData, &artifact); err != nil {
				return fmt.Errorf("parsing artifact: %w", err)
			}

			capsData, err := os.ReadFile(capsPath)
			if err != nil {
				return fmt.Errorf("reading capabilities: %w", err)
			}
			var capFile protocol.CapabilityFile
			if err := json.Unmarshal(capsData, &capFile); err != nil {
				return fmt.Errorf("parsing capabilities: %w", err)
			}

			if verbose {
				fmt.Fprintf(os.Stderr, "loaded artifact: %s (%d sections)\n", artifactPath, len(artifact.Sections))
				fmt.Fprintf(os.Stderr, "loaded capabilities: %s (%d capabilities)\n", capsPath, len(capFile.Capabilities))
			}

			cfg, err := config.Load(cfgFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: could not load config (%v), using defaults\n", err)
				cfg = config.Config{LLM: &config.LLMConfig{Provider: "llama", Model: "llama3.1:8b"}}
			}

			var llmClient *llama.Client
			if cfg.LLM != nil && cfg.LLM.Provider == "llama" {
				llmClient = llama.NewClient(cfg.LLM.Model, cfg.LLM.BaseURL)
			}

			workingDir := filepath.Dir(capsPath)
			var a protocol.Adapter
			if adapterCmd != "" {
				a = adapter.NewSubprocessAdapter(adapterCmd, workingDir)
			}

			opts := runtime.Options{
				Verbose:    verbose,
				WorkingDir: workingDir,
			}

			if verbose {
				fmt.Fprintln(os.Stderr, "resolving capability graph...")
			}

			if err := runtime.Run(context.Background(), artifact, capFile.Capabilities, llmClient, a, opts); err != nil {
				return fmt.Errorf("runtime: %w", err)
			}

			if verbose {
				fmt.Fprintln(os.Stderr, "done")
			}
			return nil
		},
	}
	cmd.Flags().String("artifact", "artifact.json", "path to context artifact")
	cmd.Flags().String("capabilities", "capabilities.json", "path to capability declarations")
	cmd.Flags().String("adapter", "", "adapter command (e.g., 'deno run adapter/src/index.ts')")
	cmd.Flags().Bool("verbose", false, "verbose output")
	return cmd
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Generate a starter contextsynth.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(cfgFile); err == nil {
				return fmt.Errorf("%s already exists — refusing to overwrite", cfgFile)
			}

			template := `version: "1"

output:
  path: artifact.json

budget: 10000

sources:
  - markdown: docs/*.md
    weight: 1.0

# sections:
#   - name: Domain Knowledge
#     budget: 0.30
#   - name: Architecture Decisions
#     budget: 0.25
#   - name: Requirements
#     budget: 0.25
#   - name: Constraints
#     budget: 0.20

# llm:
#   provider: llama
#   model: llama3.1:8b
`
			if err := os.WriteFile(cfgFile, []byte(template), 0644); err != nil {
				return fmt.Errorf("writing config: %w", err)
			}
			fmt.Printf("created %s\n", cfgFile)
			return nil
		},
	}
}

func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var parts []string
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}
	return strings.Join(parts, ",")
}
