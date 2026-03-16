package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type mockLLMClient struct {
	response string
	err      error
	prompts  []string
}

func (m *mockLLMClient) Complete(_ context.Context, prompt, _ string) (string, error) {
	m.prompts = append(m.prompts, prompt)
	return m.response, m.err
}

type mockAdapter struct {
	called bool
	state  protocol.ResolvedState
	err    error
}

func (m *mockAdapter) Project(_ context.Context, state protocol.ResolvedState) error {
	m.called = true
	m.state = state
	return m.err
}

func TestRun_LLMExecutor(t *testing.T) {
	client := &mockLLMClient{
		response: `{"tempo_bpm": 120, "key": "C"}`,
	}

	caps := []protocol.Capability{
		{
			ID:        "interpret",
			Inputs:    []string{"artifact"},
			Outputs:   []string{"params"},
			DependsOn: []string{},
			Executor: protocol.ExecutorSpec{
				Type:   "llm",
				Config: json.RawMessage(`{"prompt_template": "Analyze: {{artifact}}"}`),
			},
		},
	}

	adapter := &mockAdapter{}
	artifact := protocol.Artifact{
		Meta: protocol.Meta{Generator: "test"},
	}

	err := Run(context.Background(), artifact, caps, client, adapter, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(client.prompts) != 1 {
		t.Fatalf("expected 1 prompt, got %d", len(client.prompts))
	}

	if !adapter.called {
		t.Error("adapter should have been called")
	}

	params, ok := adapter.state.Outputs["params"]
	if !ok {
		t.Fatal("params output not found in resolved state")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		t.Fatalf("params should be a map, got %T", params)
	}
	if paramsMap["key"] != "C" {
		t.Errorf("key = %v, want C", paramsMap["key"])
	}
}

func TestRun_SubprocessExecutor(t *testing.T) {
	scriptDir := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "echo.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\necho '{\"result\": \"ok\"}'\n"), 0755); err != nil {
		t.Fatal(err)
	}

	caps := []protocol.Capability{
		{
			ID:        "process",
			Inputs:    []string{},
			Outputs:   []string{"output"},
			DependsOn: []string{},
			Executor: protocol.ExecutorSpec{
				Type:   "subprocess",
				Config: json.RawMessage(fmt.Sprintf(`{"command": "%s"}`, scriptPath)),
			},
		},
	}

	adapter := &mockAdapter{}
	artifact := protocol.Artifact{}

	err := Run(context.Background(), artifact, caps, nil, adapter, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output, ok := adapter.state.Outputs["output"]
	if !ok {
		t.Fatal("output not found in resolved state")
	}
	outputMap, ok := output.(map[string]interface{})
	if !ok {
		t.Fatalf("output should be a map, got %T", output)
	}
	if outputMap["result"] != "ok" {
		t.Errorf("result = %v, want ok", outputMap["result"])
	}
}

func TestRun_DependencyChain(t *testing.T) {
	client := &mockLLMClient{
		response: `{"mood": "melancholic"}`,
	}

	scriptDir := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "cat_input.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\ncat\n"), 0755); err != nil {
		t.Fatal(err)
	}

	caps := []protocol.Capability{
		{
			ID:        "step1",
			Inputs:    []string{"artifact"},
			Outputs:   []string{"step1_out"},
			DependsOn: []string{},
			Executor: protocol.ExecutorSpec{
				Type:   "llm",
				Config: json.RawMessage(`{"prompt_template": "Analyze: {{artifact}}"}`),
			},
		},
		{
			ID:        "step2",
			Inputs:    []string{"step1_out"},
			Outputs:   []string{"final"},
			DependsOn: []string{"step1"},
			Executor: protocol.ExecutorSpec{
				Type:   "subprocess",
				Config: json.RawMessage(fmt.Sprintf(`{"command": "%s"}`, scriptPath)),
			},
		},
	}

	adapter := &mockAdapter{}
	artifact := protocol.Artifact{}

	err := Run(context.Background(), artifact, caps, client, adapter, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := adapter.state.Outputs["step1"]; !ok {
		t.Error("step1 output should exist in resolved state")
	}
	if _, ok := adapter.state.Outputs["step1_out"]; !ok {
		t.Error("step1_out output should exist in resolved state")
	}
	if _, ok := adapter.state.Outputs["step2"]; !ok {
		t.Error("step2 output should exist in resolved state")
	}
	if _, ok := adapter.state.Outputs["final"]; !ok {
		t.Error("final output should exist in resolved state")
	}
}

func TestRun_NoLLMClient(t *testing.T) {
	caps := []protocol.Capability{
		{
			ID:      "needs_llm",
			Inputs:  []string{},
			Outputs: []string{},
			Executor: protocol.ExecutorSpec{
				Type:   "llm",
				Config: json.RawMessage(`{"prompt_template": "test"}`),
			},
		},
	}

	err := Run(context.Background(), protocol.Artifact{}, caps, nil, nil, Options{})
	if err == nil {
		t.Fatal("expected error when LLM client is nil")
	}
}

func TestRun_UnknownExecutorType(t *testing.T) {
	caps := []protocol.Capability{
		{
			ID:      "bad",
			Inputs:  []string{},
			Outputs: []string{},
			Executor: protocol.ExecutorSpec{
				Type:   "unknown",
				Config: json.RawMessage(`{}`),
			},
		},
	}

	err := Run(context.Background(), protocol.Artifact{}, caps, nil, nil, Options{})
	if err == nil {
		t.Fatal("expected error for unknown executor type")
	}
}

func TestRun_NilAdapter(t *testing.T) {
	caps := []protocol.Capability{}
	err := Run(context.Background(), protocol.Artifact{}, caps, nil, nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error with nil adapter: %v", err)
	}
}

func TestRun_AdapterError(t *testing.T) {
	adapter := &mockAdapter{err: fmt.Errorf("adapter failed")}
	caps := []protocol.Capability{}

	err := Run(context.Background(), protocol.Artifact{}, caps, nil, adapter, Options{})
	if err == nil {
		t.Fatal("expected error from adapter")
	}
}

func TestBuildInputs(t *testing.T) {
	state := protocol.ResolvedState{
		Outputs: map[string]interface{}{
			"a": "value_a",
			"b": 42,
		},
	}
	cap := protocol.Capability{
		Inputs: []string{"a", "b", "missing"},
	}

	inputs := buildInputs(cap, state)

	if inputs["a"] != "value_a" {
		t.Errorf("input a = %v, want value_a", inputs["a"])
	}
	if inputs["b"] != 42 {
		t.Errorf("input b = %v, want 42", inputs["b"])
	}
	if _, ok := inputs["missing"]; ok {
		t.Error("missing input should not be present")
	}
}

func TestRun_LLMExecutor_NonJSONResponse(t *testing.T) {
	client := &mockLLMClient{
		response: "plain text response",
	}

	caps := []protocol.Capability{
		{
			ID:      "interpret",
			Inputs:  []string{},
			Outputs: []string{"out"},
			Executor: protocol.ExecutorSpec{
				Type:   "llm",
				Config: json.RawMessage(`{"prompt_template": "test"}`),
			},
		},
	}

	adapter := &mockAdapter{}
	err := Run(context.Background(), protocol.Artifact{}, caps, client, adapter, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := adapter.state.Outputs["out"]
	if out != "plain text response" {
		t.Errorf("output = %v, want plain text response", out)
	}
}

func TestRun_LLMExecutor_CodeFencedResponse(t *testing.T) {
	client := &mockLLMClient{
		response: "```json\n{\"key\": \"value\"}\n```",
	}

	caps := []protocol.Capability{
		{
			ID:      "interpret",
			Inputs:  []string{},
			Outputs: []string{"out"},
			Executor: protocol.ExecutorSpec{
				Type:   "llm",
				Config: json.RawMessage(`{"prompt_template": "test"}`),
			},
		},
	}

	adapter := &mockAdapter{}
	err := Run(context.Background(), protocol.Artifact{}, caps, client, adapter, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, ok := adapter.state.Outputs["out"].(map[string]interface{})
	if !ok {
		t.Fatalf("output should be a map, got %T", adapter.state.Outputs["out"])
	}
	if out["key"] != "value" {
		t.Errorf("key = %v, want value", out["key"])
	}
}
