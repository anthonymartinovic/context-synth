package llm

import "context"

type Client interface {
	Complete(ctx context.Context, prompt, systemPrompt string) (string, error)
}
