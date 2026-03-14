package audio

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anthonymartinovic/context-synth/internal/config"
	"github.com/anthonymartinovic/context-synth/internal/source"
	"github.com/anthonymartinovic/context-synth/internal/token"
)

type Provider struct {
	AnalyzerPath string
}

type audioFeatures struct {
	Tempo            float64   `json:"tempo"`
	Key              int       `json:"key"`
	Mode             int       `json:"mode"`
	Energy           float64   `json:"energy"`
	Danceability     float64   `json:"danceability"`
	Valence          float64   `json:"valence"`
	Loudness         float64   `json:"loudness"`
	SpectralContrast []float64 `json:"spectral_contrast"`
	SpectralRolloff  float64   `json:"spectral_rolloff"`
	Tonnetz          []float64 `json:"tonnetz"`
	MFCCs            []float64 `json:"mfccs"`
	ZeroCrossingRate float64   `json:"zero_crossing_rate"`
	HarmonicRatio    float64   `json:"harmonic_ratio"`
}

func (p *Provider) Resolve(ctx context.Context, decl config.SourceDecl, declOrder int) ([]source.ResolvedSource, error) {
	if decl.Audio == "" {
		return nil, fmt.Errorf("audio source declaration has empty audio field")
	}

	if _, err := os.Stat(decl.Audio); err != nil {
		return nil, fmt.Errorf("audio file not found: %w", err)
	}

	features, err := p.analyze(ctx, decl.Audio)
	if err != nil {
		return nil, fmt.Errorf("analyzing %s: %w", filepath.Base(decl.Audio), err)
	}

	content := formatContent(filepath.Base(decl.Audio), features)
	hash := sha256.Sum256([]byte(content))

	return []source.ResolvedSource{{
		Path:        decl.Audio,
		Content:     []byte(content),
		ContentHash: fmt.Sprintf("%x", hash),
		TokenCount:  token.Estimate(content),
		Weight:      decl.Weight,
		DeclOrder:   declOrder,
	}}, nil
}

func (p *Provider) analyze(ctx context.Context, audioPath string) (*audioFeatures, error) {
	analyzerPath := p.AnalyzerPath
	if analyzerPath == "" {
		analyzerPath = os.Getenv("CS_AUDIO_ANALYZER")
	}
	if analyzerPath == "" {
		return nil, fmt.Errorf("no audio analyzer configured: set AnalyzerPath or CS_AUDIO_ANALYZER env var")
	}

	cmd := exec.CommandContext(ctx, "python3", analyzerPath, audioPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("audio analyzer failed: %w\nstderr: %s", err, stderr.String())
	}

	var features audioFeatures
	if err := json.Unmarshal(stdout.Bytes(), &features); err != nil {
		return nil, fmt.Errorf("parsing analyzer output: %w\nraw output: %s", err, stdout.String())
	}

	return &features, nil
}

func formatContent(filename string, f *audioFeatures) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Audio: %s\n", filename))
	b.WriteString(fmt.Sprintf("Tempo: %.1f BPM (percussive)\n", f.Tempo))
	b.WriteString(fmt.Sprintf("Key: %s %s (harmonic)\n", keyName(f.Key), modeName(f.Mode)))
	b.WriteString(fmt.Sprintf("Energy: %.2f\n", f.Energy))
	b.WriteString(fmt.Sprintf("Danceability: %.2f\n", f.Danceability))
	b.WriteString(fmt.Sprintf("Valence: %.2f\n", f.Valence))
	b.WriteString(fmt.Sprintf("Loudness: %.1f dB\n", f.Loudness))
	b.WriteString(fmt.Sprintf("Spectral Contrast: %s\n", formatFloatSlice(f.SpectralContrast)))
	b.WriteString(fmt.Sprintf("Spectral Rolloff: %.4f\n", f.SpectralRolloff))
	b.WriteString(fmt.Sprintf("Tonnetz: %s\n", formatFloatSlice(f.Tonnetz)))
	b.WriteString(fmt.Sprintf("MFCCs: %s\n", formatFloatSlice(f.MFCCs)))
	b.WriteString(fmt.Sprintf("Zero Crossing Rate: %.4f\n", f.ZeroCrossingRate))
	b.WriteString(fmt.Sprintf("Harmonic Ratio: %.2f\n", f.HarmonicRatio))
	return b.String()
}

func formatFloatSlice(vals []float64) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%.2f", v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func keyName(key int) string {
	keys := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	if key < 0 || key >= len(keys) {
		return "unknown"
	}
	return keys[key]
}

func modeName(mode int) string {
	if mode == 1 {
		return "major"
	}
	return "minor"
}
