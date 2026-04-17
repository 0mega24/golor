// Package adjust provides HSL and HSV color adjustment functions.
package adjust

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// Lighten increases the HSL lightness of c by amount (additive, clamped to [0,1]).
func Lighten(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.L = clamp01(hsl.L + amount)
	return convert.FromHSL(hsl)
}

// SetLightness sets the HSL lightness of c to l (clamped to [0,1]).
func SetLightness(c golor.Color, l float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.L = clamp01(l)
	return convert.FromHSL(hsl)
}
