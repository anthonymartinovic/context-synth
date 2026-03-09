package render

import (
	"fmt"
	"strings"

	"github.com/anthonymartinovic/context-synth/internal/model"
)

func Render(artifact model.Artifact) string {
	var b strings.Builder

	writeFrontMatter(&b, artifact.FrontMatter)

	if len(artifact.Sections) > 0 {
		for _, sec := range artifact.Sections {
			writeSection(&b, sec)
		}
	}

	if len(artifact.Omissions) > 0 {
		writeOmissions(&b, artifact.Omissions)
	}

	return b.String()
}

func writeFrontMatter(b *strings.Builder, fm model.FrontMatter) {
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("generator: %s\n", fm.Generator))
	b.WriteString(fmt.Sprintf("mode: %s\n", fm.Mode))
	b.WriteString(fmt.Sprintf("snapshot: %s\n", fm.Snapshot))
	b.WriteString(fmt.Sprintf("config: %s\n", fm.Config))
	b.WriteString(fmt.Sprintf("timestamp: %s\n", fm.Timestamp.UTC().Format("2006-01-02T15:04:05Z")))
	b.WriteString(fmt.Sprintf("budget: %d\n", fm.Budget))
	b.WriteString(fmt.Sprintf("tokens_used: %d\n", fm.TokensUsed))
	b.WriteString("---\n\n")
}

func writeSection(b *strings.Builder, sec model.ArtifactSection) {
	if sec.Name != "" {
		b.WriteString(fmt.Sprintf("# %s\n\n", sec.Name))
	}

	for i, item := range sec.Items {
		b.WriteString(strings.TrimSpace(item.Content))
		b.WriteString("\n\n")

		isLast := i == len(sec.Items)-1
		nextIsDifferentSource := isLast || sec.Items[i+1].SourcePath != item.SourcePath

		if nextIsDifferentSource {
			b.WriteString(fmt.Sprintf("> source: %s | weight: %.1f | hash: %s\n\n",
				item.SourcePath,
				item.Weight,
				item.ContentHash[:8],
			))
		}
	}
}

func writeOmissions(b *strings.Builder, omissions []model.Omission) {
	b.WriteString("# Omissions\n\n")
	b.WriteString("| Source | Reason | Weight | Tokens |\n")
	b.WriteString("|--------|--------|--------|--------|\n")
	for _, o := range omissions {
		b.WriteString(fmt.Sprintf("| %s | %s | %.1f | %d |\n",
			o.SourcePath,
			o.Reason,
			o.Weight,
			o.TokenCount,
		))
	}
}
