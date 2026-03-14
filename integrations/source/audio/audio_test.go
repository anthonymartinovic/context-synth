package audio

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/config"
)

func TestFormatContent(t *testing.T) {
	f := &audioFeatures{
		Tempo:            128.0,
		Key:              2,
		Mode:             0,
		Energy:           0.72,
		Danceability:     0.65,
		Valence:          0.45,
		Loudness:         -8.3,
		SpectralContrast: []float64{21.5, 18.3, 15.7, 12.1, 10.9, 8.4, 5.2},
		SpectralRolloff:  0.8512,
		Tonnetz:          []float64{0.02, -0.01, 0.05, -0.03, 0.01, 0.04},
		MFCCs:            []float64{-200.1, 80.5, -12.3, 30.0, -5.6, 20.1, -8.9, 15.2, -3.4, 10.0, -2.1, 8.3, -1.5},
		ZeroCrossingRate: 0.0832,
		HarmonicRatio:    0.43,
	}

	content := formatContent("midnight_drive.mp3", f)

	checks := []string{
		"Audio: midnight_drive.mp3",
		"Tempo: 128.0 BPM (percussive)",
		"Key: D minor (harmonic)",
		"Energy: 0.72",
		"Danceability: 0.65",
		"Valence: 0.45",
		"Loudness: -8.3 dB",
		"Spectral Contrast:",
		"Spectral Rolloff: 0.8512",
		"Tonnetz:",
		"MFCCs:",
		"Zero Crossing Rate: 0.0832",
		"Harmonic Ratio: 0.43",
	}

	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("content missing %q\ngot:\n%s", check, content)
		}
	}
}

func TestFormatFloatSlice(t *testing.T) {
	result := formatFloatSlice([]float64{1.234, 5.678, 9.0})
	expected := "[1.23, 5.68, 9.00]"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestFormatFloatSliceEmpty(t *testing.T) {
	result := formatFloatSlice(nil)
	expected := "[]"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestKeyName(t *testing.T) {
	tests := []struct {
		key  int
		want string
	}{
		{0, "C"}, {1, "C#"}, {2, "D"}, {11, "B"},
		{-1, "unknown"}, {12, "unknown"},
	}
	for _, tt := range tests {
		if got := keyName(tt.key); got != tt.want {
			t.Errorf("keyName(%d) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestModeName(t *testing.T) {
	if modeName(1) != "major" {
		t.Error("modeName(1) should be major")
	}
	if modeName(0) != "minor" {
		t.Error("modeName(0) should be minor")
	}
}

func TestResolveEmptyMP3(t *testing.T) {
	p := &Provider{}
	_, err := p.Resolve(context.Background(), config.SourceDecl{Weight: 0.8}, 0)
	if err == nil {
		t.Fatal("expected error for empty audio field")
	}
}

func TestResolveMissingFile(t *testing.T) {
	p := &Provider{}
	decl := config.SourceDecl{Audio: "/nonexistent/path/song.mp3", Weight: 0.8}
	_, err := p.Resolve(context.Background(), decl, 0)
	if err == nil {
		t.Fatal("expected error for missing audio file")
	}
	if !strings.Contains(err.Error(), "audio file not found") {
		t.Errorf("error = %v, want 'audio file not found'", err)
	}
}

func TestResolveMissingAnalyzer(t *testing.T) {
	dir := t.TempDir()
	mp3File := filepath.Join(dir, "test.mp3")
	os.WriteFile(mp3File, []byte("fake audio"), 0644)

	t.Setenv("CS_AUDIO_ANALYZER", "")
	p := &Provider{}
	decl := config.SourceDecl{Audio: mp3File, Weight: 0.8}
	_, err := p.Resolve(context.Background(), decl, 0)
	if err == nil {
		t.Fatal("expected error for missing analyzer")
	}
	if !strings.Contains(err.Error(), "no audio analyzer configured") {
		t.Errorf("error = %v, want 'no audio analyzer configured'", err)
	}
}

func TestResolveWithMockAnalyzer(t *testing.T) {
	dir := t.TempDir()

	mp3File := filepath.Join(dir, "test_track.mp3")
	os.WriteFile(mp3File, []byte("fake audio"), 0644)

	script := filepath.Join(dir, "mock_analyze.py")
	os.WriteFile(script, []byte(`import json, sys
print(json.dumps({
    "tempo": 120.0, "key": 7, "mode": 1,
    "energy": 0.65, "danceability": 0.70, "valence": 0.55, "loudness": -7.2,
    "spectral_contrast": [20.0, 17.5, 14.2, 11.8, 9.5, 7.1, 4.3],
    "spectral_rolloff": 0.82,
    "tonnetz": [0.01, -0.02, 0.03, -0.01, 0.02, 0.01],
    "mfccs": [-180.0, 75.0, -10.0, 28.0, -4.5, 18.0, -7.0, 13.0, -2.5, 9.0, -1.8, 7.0, -1.2],
    "zero_crossing_rate": 0.075,
    "harmonic_ratio": 0.58
}))
`), 0644)

	p := &Provider{AnalyzerPath: script}
	decl := config.SourceDecl{Audio: mp3File, Weight: 0.8}
	sources, err := p.Resolve(context.Background(), decl, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sources) != 1 {
		t.Fatalf("got %d sources, want 1", len(sources))
	}

	s := sources[0]
	if s.Path != mp3File {
		t.Errorf("path = %q, want %q", s.Path, mp3File)
	}
	if s.Weight != 0.8 {
		t.Errorf("weight = %f, want 0.8", s.Weight)
	}
	if s.DeclOrder != 3 {
		t.Errorf("declOrder = %d, want 3", s.DeclOrder)
	}
	if s.ContentHash == "" {
		t.Error("content hash is empty")
	}
	if s.TokenCount == 0 {
		t.Error("token count is zero")
	}

	content := string(s.Content)
	if !strings.Contains(content, "Audio: test_track.mp3") {
		t.Errorf("content missing filename\n%s", content)
	}
	if !strings.Contains(content, "Key: G major") {
		t.Errorf("content missing key\n%s", content)
	}
	if !strings.Contains(content, "Tempo: 120.0 BPM") {
		t.Errorf("content missing tempo\n%s", content)
	}
	if !strings.Contains(content, "Harmonic Ratio: 0.58") {
		t.Errorf("content missing harmonic ratio\n%s", content)
	}
}
