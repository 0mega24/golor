package harmony

import (
	"math"

	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// rotateHue returns c with its HSL hue rotated by degrees and alpha preserved.
func rotateHue(c golor.Color, degrees float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.H = math.Mod(hsl.H+degrees+360*10, 360)
	out := convert.FromHSL(hsl)
	out.A = c.A
	return out
}
