package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GeminiClient struct {
	Model  string
	APIKey string
}

func NewGeminiClient(model, apiKeyEnv string) (*GeminiClient, error) {
	key := os.Getenv(apiKeyEnv)
	if key == "" {
		return nil, fmt.Errorf("environment variable %s is not set", apiKeyEnv)
	}
	if model == "" {
		model = "gemini-2.5-pro"
	}
	return &GeminiClient{Model: model, APIKey: key}, nil
}

func (g *GeminiClient) Complete(ctx context.Context, prompt, systemPrompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		g.Model, g.APIKey)

	var contents []geminiContent
	if systemPrompt != "" {
		contents = append(contents, geminiContent{
			Role:  "user",
			Parts: []geminiPart{{Text: systemPrompt}},
		})
		contents = append(contents, geminiContent{
			Role:  "model",
			Parts: []geminiPart{{Text: "Understood. I will follow these instructions."}},
		})
	}
	contents = append(contents, geminiContent{
		Role:  "user",
		Parts: []geminiPart{{Text: prompt}},
	})

	reqBody := geminiRequest{Contents: contents}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini API call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}
