package adjust

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Darken decreases the HSL lightness of c by amount (additive, clamped to [0,1]) and preserves alpha.
func Darken(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.L = clamp01(hsl.L - amount)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}
