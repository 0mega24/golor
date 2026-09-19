package blend

import "github.com/0mega24/golor/v2"

// ColorBurn applies the Color Burn blend mode with Porter-Duff source-over alpha compositing.
func ColorBurn(base, layer golor.Color) golor.Color {
	return compositeBlend(base, layer, golor.RGBf(
		colorBurnChannel(base.R, layer.R),
		colorBurnChannel(base.G, layer.G),
		colorBurnChannel(base.B, layer.B),
	))
}

func colorBurnChannel(base, layer float64) float64 {
	if base >= 1 {
		return 1
	}
	if layer <= 0 {
		return 0
	}
	return clamp01(1 - (1-base)/layer)
}
