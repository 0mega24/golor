package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Luminosity applies the Luminosity blend mode with Porter-Duff source-over alpha compositing.
// Takes hue and saturation from base, lightness from layer (operates in HSL space).
func Luminosity(base, layer golor.Color) golor.Color {
	hsl := convert.ToHSL(base)
	hsl.L = convert.ToHSL(layer).L
	return compositeBlend(base, layer, convert.FromHSL(hsl))
}
