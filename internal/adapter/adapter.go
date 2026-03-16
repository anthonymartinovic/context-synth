package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type SubprocessAdapter struct {
	Command    string
	Args       []string
	WorkingDir string
}

func NewSubprocessAdapter(adapterCmd string, workingDir string) *SubprocessAdapter {
	parts := strings.Fields(adapterCmd)
	if len(parts) == 0 {
		return &SubprocessAdapter{}
	}
	return &SubprocessAdapter{
		Command:    parts[0],
		Args:       parts[1:],
		WorkingDir: workingDir,
	}
}

func (a *SubprocessAdapter) Project(ctx context.Context, state protocol.ResolvedState) error {
	if a.Command == "" {
		return fmt.Errorf("adapter has no command configured")
	}

	stateJSON, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshaling resolved state: %w", err)
	}

	cmd := exec.CommandContext(ctx, a.Command, a.Args...)
	cmd.Stdin = bytes.NewReader(stateJSON)
	cmd.Stdout = os.Stdout

	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if a.WorkingDir != "" {
		cmd.Dir = a.WorkingDir
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("adapter process failed: %w\nstderr: %s", err, stderrBuf.String())
	}

	return nil
}
