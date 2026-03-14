package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type subprocessExecutorConfig struct {
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"working_dir,omitempty"`
}

func executeSubprocess(ctx context.Context, cap protocol.Capability, inputs map[string]interface{}, workingDir string) (interface{}, error) {
	var cfg subprocessExecutorConfig
	if err := json.Unmarshal(cap.Executor.Config, &cfg); err != nil {
		return nil, fmt.Errorf("parsing subprocess executor config for %q: %w", cap.ID, err)
	}

	if cfg.Command == "" {
		return nil, fmt.Errorf("subprocess executor for %q has empty command", cap.ID)
	}

	inputJSON, err := json.Marshal(inputs)
	if err != nil {
		return nil, fmt.Errorf("marshaling inputs for %q: %w", cap.ID, err)
	}

	cmd := exec.CommandContext(ctx, cfg.Command, cfg.Args...)
	cmd.Stdin = bytes.NewReader(inputJSON)

	dir := workingDir
	if cfg.WorkingDir != "" {
		dir = cfg.WorkingDir
	}
	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("subprocess %q failed: %w\nstderr: %s", cap.ID, err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return nil, nil
	}

	var result interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		fmt.Fprintf(os.Stderr, "warning: subprocess output for %q is not valid JSON, using raw string\n", cap.ID)
		return output, nil
	}

	return result, nil
}
