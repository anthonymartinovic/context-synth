package render

import (
	"encoding/json"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

func Render(artifact protocol.Artifact) (string, error) {
	data, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
