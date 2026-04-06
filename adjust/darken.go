package adjust

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// Darken decreases the HSL lightness of c by amount (additive, clamped to [0,1]).
func Darken(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.L = clamp01(hsl.L - amount)
	return convert.FromHSL(hsl)
}
