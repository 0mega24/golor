package adjust

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// ShiftHue rotates the hue of c by degrees (circular, wraps around 360°).
func ShiftHue(c golor.Color, degrees float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.H = math.Mod(hsl.H+degrees+360*10, 360)
	return convert.FromHSL(hsl)
}

// SetHue sets the hue of c to h degrees (clamped/wrapped to [0, 360)).
func SetHue(c golor.Color, h float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.H = math.Mod(h+360*10, 360)
	return convert.FromHSL(hsl)
}
