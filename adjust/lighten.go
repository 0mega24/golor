// Package adjust provides HSL, HSV, RGB, and OKLCH color adjustment functions.
package adjust

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
	"github.com/0mega24/golor/v2/internal/space"
)

// Lighten increases the HSL lightness of c by amount (additive, clamped to [0,1]) and preserves alpha.
func Lighten(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.L = space.Clamp01(hsl.L + amount)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}

// SetLightness sets the HSL lightness of c to l (clamped to [0,1]) and preserves alpha.
func SetLightness(c golor.Color, l float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.L = space.Clamp01(l)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}

// Darken decreases the HSL lightness of c by amount (additive, clamped to [0,1]) and preserves alpha.
func Darken(c golor.Color, amount float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.L = space.Clamp01(hsl.L - amount)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}
