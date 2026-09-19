package adjust

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Desaturate decreases the HSL saturation of c by amount (additive, clamped to [0,1]) and preserves alpha.
func Desaturate(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.S = clamp01(hsl.S - amount)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}
