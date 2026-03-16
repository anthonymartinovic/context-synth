package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthonymartinovic/context-synth/internal/graph"
	"github.com/anthonymartinovic/context-synth/internal/llm"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type Options struct {
	Verbose    bool
	WorkingDir string
}

func Run(ctx context.Context, artifact protocol.Artifact, capabilities []protocol.Capability, client llm.Client, adapter protocol.Adapter, opts Options) error {
	g, err := graph.Build(capabilities)
	if err != nil {
		return fmt.Errorf("building capability graph: %w", err)
	}

	state := protocol.ResolvedState{
		Artifact: artifact,
		Outputs:  make(map[string]interface{}),
	}

	artifactJSON, err := json.Marshal(artifact)
	if err != nil {
		return fmt.Errorf("marshaling artifact: %w", err)
	}
	state.Outputs["artifact"] = json.RawMessage(artifactJSON)

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "graph order: %v\n", g.Order())
	}

	for _, capID := range g.Order() {
		cap, _ := g.Get(capID)

		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "  resolving: %s (executor: %s)\n", capID, cap.Executor.Type)
		}

		inputs := buildInputs(cap, state)

		var output interface{}
		switch cap.Executor.Type {
		case "llm":
			if client == nil {
				return fmt.Errorf("capability %q requires LLM executor but no LLM client provided", capID)
			}
			output, err = executeLLM(ctx, cap, client, inputs)
		case "subprocess":
			output, err = executeSubprocess(ctx, cap, inputs, opts.WorkingDir)
		default:
			return fmt.Errorf("unknown executor type %q for capability %q", cap.Executor.Type, capID)
		}

		if err != nil {
			return fmt.Errorf("resolving capability %q: %w", capID, err)
		}

		state.Outputs[capID] = output

		for _, outputName := range cap.Outputs {
			state.Outputs[outputName] = output
		}

		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "  resolved: %s\n", capID)
			if outputJSON, err := json.MarshalIndent(output, "    ", "  "); err == nil {
				const maxLen = 500
				s := string(outputJSON)
				if len(s) > maxLen {
					s = s[:maxLen] + "..."
				}
				fmt.Fprintf(os.Stderr, "    output: %s\n", s)
			}
		}
	}

	if adapter != nil {
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "invoking adapter...\n")
		}
		if err := adapter.Project(ctx, state); err != nil {
			return fmt.Errorf("projection: %w", err)
		}
	}

	return nil
}

func buildInputs(cap protocol.Capability, state protocol.ResolvedState) map[string]interface{} {
	inputs := make(map[string]interface{})
	for _, inputName := range cap.Inputs {
		if val, ok := state.Outputs[inputName]; ok {
			inputs[inputName] = val
		}
	}
	return inputs
}
