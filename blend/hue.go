package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Hue applies the Hue blend mode with Porter-Duff source-over alpha compositing.
func Hue(base, layer golor.Color) golor.Color {
	hsl := convert.ToHSL(base)
	hsl.H = convert.ToHSL(layer).H
	return compositeBlend(base, layer, convert.FromHSL(hsl))
}
