package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anthonymartinovic/context-synth/internal/llm"
	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type llmExecutorConfig struct {
	PromptTemplate string `json:"prompt_template"`
}

func executeLLM(ctx context.Context, cap protocol.Capability, client llm.Client, inputs map[string]interface{}) (interface{}, error) {
	var cfg llmExecutorConfig
	if err := json.Unmarshal(cap.Executor.Config, &cfg); err != nil {
		return nil, fmt.Errorf("parsing LLM executor config for %q: %w", cap.ID, err)
	}

	if cfg.PromptTemplate == "" {
		return nil, fmt.Errorf("LLM executor for %q has empty prompt_template", cap.ID)
	}

	prompt := cfg.PromptTemplate
	for key, val := range inputs {
		jsonVal, err := json.Marshal(val)
		if err != nil {
			return nil, fmt.Errorf("marshaling input %q for capability %q: %w", key, cap.ID, err)
		}
		prompt = strings.ReplaceAll(prompt, "{{"+key+"}}", string(jsonVal))
	}

	response, err := client.Complete(ctx, prompt, "")
	if err != nil {
		return nil, fmt.Errorf("LLM execution for %q: %w", cap.ID, err)
	}

	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var result interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		fmt.Fprintf(os.Stderr, "warning: LLM output for %q is not valid JSON, using raw string\n", cap.ID)
		return response, nil
	}

	return result, nil
}
