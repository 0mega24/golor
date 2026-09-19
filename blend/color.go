package blend

import (
	"github.com/0mega24/golor/v2"
	"github.com/0mega24/golor/v2/convert"
)

// Color applies the Color blend mode with Porter-Duff source-over alpha compositing.
// It is named for the standard blend mode and is used as blend.Color(base, layer),
// distinct from the golor.Color type.
func Color(base, layer golor.Color) golor.Color {
	hsl := convert.ToHSL(layer)
	hsl.L = convert.ToHSL(base).L
	return compositeBlend(base, layer, convert.FromHSL(hsl))
}
