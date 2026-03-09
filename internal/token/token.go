package token

import (
	"math"
	"strings"
)

func Estimate(text string) int {
	words := len(strings.Fields(text))
	return int(math.Ceil(float64(words) * 1.33))
}
