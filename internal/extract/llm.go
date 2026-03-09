package extract

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/llm"
	"github.com/anthonymartinovic/context-synth/internal/model"
	"github.com/anthonymartinovic/context-synth/internal/token"
)

type LLMExtractor struct {
	Client llm.Client
}

func (l *LLMExtractor) Extract(ctx context.Context, snap model.Snapshot, sections []config.SectionDecl) ([]model.Extraction, error) {
	if len(sections) == 0 {
		return nil, fmt.Errorf("LLM extractor requires at least one section definition")
	}

	var allExtractions []model.Extraction

	for _, item := range snap.Items {
		prompt := buildExtractionPrompt(item, sections)
		systemPrompt := buildSystemPrompt(sections)

		response, err := l.Client.Complete(ctx, prompt, systemPrompt)
		if err != nil {
			return nil, fmt.Errorf("LLM extraction for %s: %w", item.Source.Path, err)
		}

		extractions, err := parseExtractionResponse(response, item)
		if err != nil {
			return nil, fmt.Errorf("parsing LLM response for %s: %w", item.Source.Path, err)
		}

		allExtractions = append(allExtractions, extractions...)
	}

	return allExtractions, nil
}

func buildSystemPrompt(sections []config.SectionDecl) string {
	var sectionNames []string
	for _, s := range sections {
		sectionNames = append(sectionNames, s.Name)
	}

	return fmt.Sprintf(`You are a content extraction assistant. Your job is to decompose source documents into discrete knowledge items and classify each into one of these sections: %s.

Respond ONLY with a JSON array. Each element must have:
- "section": one of the section names listed above (exact match)
- "content": the extracted text, preserving the key information

Do not add commentary. Do not wrap in markdown code fences. Return only the JSON array.`, strings.Join(sectionNames, ", "))
}

func buildExtractionPrompt(item model.SnapshotItem, sections []config.SectionDecl) string {
	var sectionList strings.Builder
	for _, s := range sections {
		sectionList.WriteString(fmt.Sprintf("- %s\n", s.Name))
	}

	return fmt.Sprintf(`Extract discrete knowledge items from the following source document and classify each into one of these sections:

%s
Source file: %s

---
%s
---

Return a JSON array of objects with "section" and "content" fields.`,
		sectionList.String(),
		item.Source.Path,
		string(item.Source.Content),
	)
}

type extractionItem struct {
	Section string `json:"section"`
	Content string `json:"content"`
}

func parseExtractionResponse(response string, item model.SnapshotItem) ([]model.Extraction, error) {
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var items []extractionItem
	if err := json.Unmarshal([]byte(response), &items); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w\nraw response: %s", err, response)
	}

	var extractions []model.Extraction
	for _, ei := range items {
		extractions = append(extractions, model.Extraction{
			Content:     ei.Content,
			Section:     ei.Section,
			TokenCount:  token.Estimate(ei.Content),
			Weight:      item.Source.Weight,
			SourcePath:  item.Source.Path,
			ContentHash: item.Source.ContentHash,
			DeclOrder:   item.Source.DeclOrder,
		})
	}

	return extractions, nil
}
