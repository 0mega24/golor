package blend

import "github.com/0mega24/golor/v2"

// Multiply applies the Multiply blend mode with Porter-Duff source-over alpha compositing.
func Multiply(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		base.R*layer.R,
		base.G*layer.G,
		base.B*layer.B,
	))
}
