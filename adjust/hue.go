package adjust

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// ShiftHue rotates the hue of c by degrees (circular, wraps around 360°) and preserves alpha.
func ShiftHue(c golor.Color, degrees float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.H = math.Mod(hsl.H+degrees+360*10, 360)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}

// SetHue sets the hue of c to h degrees (clamped/wrapped to [0, 360)) and preserves alpha.
func SetHue(c golor.Color, h float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.H = math.Mod(h+360*10, 360)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}
