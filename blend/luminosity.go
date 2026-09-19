package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Luminosity applies the Photoshop Luminosity blend mode.
// Takes hue and saturation from base, lightness from layer (operates in HSL space).
// The result alpha is the source-over composited alpha of base and layer.
func Luminosity(base, layer golor.Color) golor.Color {
	hsl := convert.ToHSL(base)
	hsl.L = convert.ToHSL(layer).L
	out := convert.FromHSL(hsl)
	out.A = blendedAlpha(base, layer)
	return out
}
