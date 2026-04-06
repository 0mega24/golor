package adjust

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// Saturate increases the HSL saturation of c by amount (additive, clamped to [0,1]).
func Saturate(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.S = clamp01(hsl.S + amount)
	return convert.FromHSL(hsl)
}

// SetSaturation sets the HSL saturation of c to s (clamped to [0,1]).
func SetSaturation(c golor.Color, s float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.S = clamp01(s)
	return convert.FromHSL(hsl)
}
