package blend

import "github.com/0mega24/golor/v2"

// Exclusion applies the Exclusion blend mode with Porter-Duff source-over alpha compositing.
func Exclusion(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		exclusionChannel(base.R, layer.R),
		exclusionChannel(base.G, layer.G),
		exclusionChannel(base.B, layer.B),
	))
}

func exclusionChannel(base, layer float64) float64 {
	return base + layer - 2*base*layer
}
