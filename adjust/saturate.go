package adjust

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Saturate increases the HSL saturation of c by amount (additive, clamped to [0,1]) and preserves alpha.
func Saturate(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.S = clamp01(hsl.S + amount)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}

// SetSaturation sets the HSL saturation of c to s (clamped to [0,1]) and preserves alpha.
func SetSaturation(c golor.Color, s float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.S = clamp01(s)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}
