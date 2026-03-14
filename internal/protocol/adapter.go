package protocol

import "context"

type ResolvedState struct {
	Artifact Artifact               `json:"artifact"`
	Outputs  map[string]interface{} `json:"outputs"`
}

type Adapter interface {
	Project(ctx context.Context, state ResolvedState) error
}
