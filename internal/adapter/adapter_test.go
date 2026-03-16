package adapter

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

func TestNewSubprocessAdapter_ParsesCommand(t *testing.T) {
	a := NewSubprocessAdapter("deno run --allow-read adapter/src/index.ts", "/work")

	if a.Command != "deno" {
		t.Errorf("command = %q, want deno", a.Command)
	}
	expected := []string{"run", "--allow-read", "adapter/src/index.ts"}
	if len(a.Args) != len(expected) {
		t.Fatalf("args length = %d, want %d", len(a.Args), len(expected))
	}
	for i, arg := range expected {
		if a.Args[i] != arg {
			t.Errorf("args[%d] = %q, want %q", i, a.Args[i], arg)
		}
	}
	if a.WorkingDir != "/work" {
		t.Errorf("working dir = %q, want /work", a.WorkingDir)
	}
}

func TestNewSubprocessAdapter_EmptyCommand(t *testing.T) {
	a := NewSubprocessAdapter("", "")
	if a.Command != "" {
		t.Errorf("command should be empty, got %q", a.Command)
	}
}

func TestNewSubprocessAdapter_SingleCommand(t *testing.T) {
	a := NewSubprocessAdapter("echo", "")
	if a.Command != "echo" {
		t.Errorf("command = %q, want echo", a.Command)
	}
	if len(a.Args) != 0 {
		t.Errorf("args should be empty, got %v", a.Args)
	}
}

func TestProject_InvokesSubprocess(t *testing.T) {
	scriptDir := t.TempDir()
	outputFile := filepath.Join(scriptDir, "received.json")
	scriptPath := filepath.Join(scriptDir, "capture.sh")

	script := "#!/bin/sh\ncat > " + outputFile + "\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	a := NewSubprocessAdapter(scriptPath, "")
	state := protocol.ResolvedState{
		Artifact: protocol.Artifact{
			Meta: protocol.Meta{Generator: "test"},
		},
		Outputs: map[string]interface{}{
			"key": "value",
		},
	}

	err := a.Project(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("reading captured output: %v", err)
	}

	var received protocol.ResolvedState
	if err := json.Unmarshal(data, &received); err != nil {
		t.Fatalf("parsing captured JSON: %v", err)
	}

	if received.Artifact.Meta.Generator != "test" {
		t.Errorf("generator = %q, want test", received.Artifact.Meta.Generator)
	}
}

func TestProject_EmptyCommand(t *testing.T) {
	a := NewSubprocessAdapter("", "")
	err := a.Project(context.Background(), protocol.ResolvedState{})
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestProject_SubprocessFailure(t *testing.T) {
	scriptDir := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "fail.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\nexit 1\n"), 0755); err != nil {
		t.Fatal(err)
	}

	a := NewSubprocessAdapter(scriptPath, "")
	state := protocol.ResolvedState{
		Outputs: map[string]interface{}{},
	}

	err := a.Project(context.Background(), state)
	if err == nil {
		t.Fatal("expected error for failed subprocess")
	}
}
