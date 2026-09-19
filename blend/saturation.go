package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Saturation applies the Saturation blend mode with Porter-Duff source-over alpha compositing.
func Saturation(base, layer golor.Color) golor.Color {
	hsl := convert.ToHSL(base)
	hsl.S = convert.ToHSL(layer).S
	return compositeBlend(base, layer, convert.FromHSL(hsl))
}
