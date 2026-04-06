package blend

import (
	"github.com/0mega24/golor"
	"github.com/0mega24/golor/convert"
)

// Luminosity applies the Photoshop Luminosity blend mode.
// Takes hue and saturation from base, lightness from layer (operates in HSL space).
func Luminosity(base, layer golor.Color) golor.Color {
	hsl := convert.ToHSL(base)
	hsl.L = convert.ToHSL(layer).L
	return convert.FromHSL(hsl)
}
