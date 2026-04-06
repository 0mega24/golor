package harmony

import (
	"math"

	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// rotateHue returns c with its HSL hue rotated by degrees.
func rotateHue(c golor.Color, degrees float64) golor.Color {
	hsl := convert.ToHSL(c)
	hsl.H = math.Mod(hsl.H+degrees+360*10, 360)
	return convert.FromHSL(hsl)
}
