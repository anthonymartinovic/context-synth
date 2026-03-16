package protocol

import "encoding/json"

type Capability struct {
	ID          string       `json:"id"`
	Description string       `json:"description"`
	Inputs      []string     `json:"inputs"`
	Outputs     []string     `json:"outputs"`
	DependsOn   []string     `json:"depends_on"`
	Executor    ExecutorSpec `json:"executor"`
}

type ExecutorSpec struct {
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

type CapabilityFile struct {
	Capabilities []Capability `json:"capabilities"`
}
