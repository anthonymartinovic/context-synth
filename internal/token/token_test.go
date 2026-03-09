package token

import "testing"

func TestEstimate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", "", 0},
		{"single word", "hello", 2},
		{"three words", "one two three", 4},
		{"ten words", "a b c d e f g h i j", 14},
		{"real sentence", "The quick brown fox jumps over the lazy dog", 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Estimate(tt.input)
			if got != tt.want {
				t.Errorf("Estimate(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
