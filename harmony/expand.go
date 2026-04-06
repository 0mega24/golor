package harmony

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// Expand returns n colors ranging from a dark shade through c to a light tint,
// evenly distributed in HSL lightness. The hue and saturation of c are preserved.
func Expand(c golor.Color, n int) []golor.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []golor.Color{c}
	}
	hsl := convert.ToHSL(c)
	minL := math.Max(0, hsl.L-0.4)
	maxL := math.Min(1, hsl.L+0.4)
	result := make([]golor.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		step := hsl
		step.L = minL + t*(maxL-minL)
		result[i] = convert.FromHSL(step)
	}
	return result
}
